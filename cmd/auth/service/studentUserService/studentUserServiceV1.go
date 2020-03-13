package studentUserService

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"rxt/cmd/auth/model"
	auth "rxt/cmd/auth/proto/student"
	"rxt/cmd/auth/service/baseService"
	"rxt/internal/config"
	"rxt/internal/jwt"
	"rxt/internal/util"
	"rxt/internal/wrong"
	"strconv"
	"time"
)

type V1 struct {
	baseService.Service
}

func (c *V1) Login(param *Param) (*auth.LogicResponse, error) {
	var result auth.LogicResponse
	var studentLogic model.StudentUserLogin
	c.Db.Where("student_user_login_name = ? ", param.StudentUserLoginName).First(&studentLogic)
	password := []byte(param.Password)
	dataPassword := []byte(studentLogic.Password)
	//密码验证
	err := bcrypt.CompareHashAndPassword(dataPassword, password)
	if err != nil {
		return nil, wrong.New(http.StatusExpectationFailed, errors.New("密码错误或账号已经被禁用！"))
	}
	var studentUser model.StudentUser
	c.Db.Where("student_user_no = ? ", studentLogic.StudentUserNo).
		Preload("StudentUserCampusOne").First(&studentUser)
	studentUser.LastLoginIp = int64(ip2long(param.Ip))
	timeNow := time.Now()
	studentUser.LastLoginTime = &timeNow
	c.Db.Save(&studentUser)
	token, err := jwt.GenerateToken(jwt.User{
		Sub: studentUser.StudentUserNo,
		Prv: "App\\Modes\\StudentUser",
	})
	if err != nil {
		return nil, wrong.New(http.StatusExpectationFailed, err, "token生成失败!")
	}
	courseItemStudent := model.CourseItemStudent{}
	c.Db.Joins("INNER JOIN rxt_course_item ON rxt_course_item.course_item_id = rxt_course_item_student.course_item_id").
		Where("rxt_course_item_student.student_user_campus_id = ?", studentUser.StudentUserCampusOne.StudentUserCampusId).
		Where("rxt_course_item.course_item_status = 4").
		Order("rxt_course_item.course_start_at desc").First(&courseItemStudent)
	result.Token = &auth.Token{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   0,}
	result.StudentUserHead = studentUser.StudentUserHead
	result.Student = &auth.Student{
		StudentRealName: studentUser.StudentUserCampusOne.StudentRealName,
		StudentSex:      int64(studentUser.StudentUserCampusOne.StudentSex),
		StudentUserNo:   studentUser.StudentUserNo,
		InClass:         courseItemStudent.CourseItemId,
	}
	c.SetStudentUserCache(studentUser.StudentUserNo)
	return &result, nil
}

func ip2long(ip net.IP) uint32 {
	if ip == nil {
		return 0
	}
	a := uint32(ip[12])
	b := uint32(ip[13])
	c := uint32(ip[14])
	d := uint32(ip[15])
	return uint32(a<<24 | b<<16 | c<<8 | d)
}

func long2ip(ip uint32) net.IP {
	a := byte((ip >> 24) & 0xFF)
	b := byte((ip >> 16) & 0xFF)
	c := byte((ip >> 8) & 0xFF)
	d := byte(ip & 0xFF)
	return net.IPv4(a, b, c, d)
}

type StudentUserCache struct {
	StudentUserNo int64 `json:"student_user_no"`
}

func (c V1) SetStudentUserCache(studentUserNo int64) {
	studentUserNoStr := strconv.FormatInt(studentUserNo, 10)
	StudentUserNoMd5 := util.MD5HashString(studentUserNoStr)
	c.Cache.SetJSON("student:user:"+StudentUserNoMd5, &StudentUserCache{
		StudentUserNo: studentUserNo,
	}, config.GetJwtCnf().Ttl)
}

func (c V1) GetStudentUserCache(studentUserNo int64) *StudentUserCache {
	var studentUserCache *StudentUserCache
	studentUserNoStr := strconv.FormatInt(studentUserNo, 10)
	StudentUserNoMd5 := util.MD5HashString(studentUserNoStr)
	c.Cache.Get(StudentUserNoMd5, &studentUserCache)
	return studentUserCache
}

func (c V1) DeleteStudentUserCache(studentUserNo int64) {
	studentUserNoStr := strconv.FormatInt(studentUserNo, 10)
	StudentUserNoMd5 := util.MD5HashString(studentUserNoStr)
	c.Cache.Del(StudentUserNoMd5)
}
