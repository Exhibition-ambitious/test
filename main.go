package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model           // 包含ID、CreatedAt、UpdatedAt、DeletedAt字段S
	Name       string    `gorm:"not null" json:"name"`
	Email      string    `gorm:"not null" json:"email"`
	Patients   []Patient `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
type Patient struct {
	gorm.Model          // 包含ID、CreatedAt、UpdatedAt、DeletedAt字段
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null" json:"email"`
	Status       string `json:"status"`
	RoomID       uint   `json:"roomID"`
	DoctorID     uint   `json:"doctorID"`
	Hospital_fee uint   `json:"hospital_fee"`
}
type Room struct {
	gorm.Model            // 包含ID、CreatedAt、UpdatedAt、DeletedAt字段
	Patient_num uint      `gorm:"not null" json:"patient_num"`
	Patients    []Patient `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
type Medicine struct {
	gorm.Model
	Medicine_name string `gorm:"not null" json:"medicine_name"`
	Price         uint   `gorm:"not null" json:"price"`
}

func main() {
	//创建一个服务
	ginServer := gin.Default()
	//连接数据库
	dsn := "root:1234@tcp(127.0.0.1:3306)/hospital-system?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	sqlDB, err := db.DB()
	//空闲连接池的最大数量
	sqlDB.SetMaxIdleConns(10)
	// 连接数据库最大数量
	sqlDB.SetMaxOpenConns(100)
	//连接数据库最长时间
	sqlDB.SetConnMaxLifetime(10 * time.Second)
	db.AutoMigrate(&Room{}, &Doctor{}, &Patient{}, &Medicine{})
	ginServer.POST("/add", func(ctx *gin.Context) { //增加患者
		var data Patient
		err := ctx.ShouldBindJSON(&data)
		if err != nil {
			ctx.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": gin.H{},
				"code": 400,
			})
		} else {
			//同步到数据库中
			/* 	func (u *Patient) BeforeCreate(tx *gorm.DB) (err error) {
			u.UUID = uuid.New()

				if u.Role == "admin" {
					return errors.New("invalid role")
				}
				return
			}*/
			db.Create(&data)
			ctx.JSON(200, gin.H{
				"msg":  "添加成功",
				"data": data,
				"code": 200,
			})
		}
	})
	ginServer.PUT("/update/:id", func(ctx *gin.Context) { //更新患者信息（住院信息等）
		var data Patient
		//寻找对应id的患者
		id := ctx.Param("id")
		db.Select("id").Where("id=?", id).Find(&data)
		if id == "0" {
			ctx.JSON(200, gin.H{
				"msg":  "用户信息没找到",
				"code": 400,
			})
		} else {
			err := ctx.ShouldBindJSON(&data)
			if err != nil {
				ctx.JSON(200, gin.H{
					"msg":  "修改失败",
					"code": 400,
				})
			} else {
				db.Where("id=?", id).Updates(&data)
				ctx.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}
		}
	})
	//服务器端口
	ginServer.Run(":8085")
}
