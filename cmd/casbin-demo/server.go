package main

import (
	"fmt"

	"github.com/andrepinto/casbin-demo/pkg/conf"
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/model"
	"github.com/casbin/xorm-adapter"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var enforcer *casbin.Enforcer

func LoadCasbin() {

	modelConfig := conf.MustAsset("config/rbac_model.conf")
	m := model.Model{}
	m.LoadModelFromText(string(modelConfig))

	adapter := xormadapter.NewAdapter("postgres",
		" host=mrs-dev.c9zb3qf0ltrn.eu-west-1.rds.amazonaws.com user=tlantic dbname=authorization sslmode=disable password=tlanticdev2!",
		true)

	enforcer = casbin.NewEnforcer("config/rbac_model.conf", adapter)

	enforcer.AddPolicy("alice", "data1", "read", "allow")
	enforcer.AddPolicy("admin", "cockpit", "read", "allow")
	enforcer.AddPolicy("admin", "cockpit", "write", "allow")
	enforcer.AddPolicy("admin", "instore", "read", "allow")
	enforcer.AddPolicy("admin", "instore", "write", "allow")
	enforcer.AddPolicy("cockpit-user", "cockpit", "write", "allow")
	enforcer.AddPolicy("cockpit-user", "cockpit", "read", "allow")
	//enforcer.SavePolicy()
	enforcer.AddGroupingPolicy("alice", "admin")
	enforcer.AddGroupingPolicy("bob", "cockpit-user")
	enforcer.SavePolicy()

	fmt.Printf("get policy: %v\n", enforcer.GetPolicy())

	fmt.Printf("get policy: %v\n", enforcer.GetPolicy())

	enforcer.LoadPolicy()
}

func main() {
	LoadCasbin()

	r := gin.Default()
	r.GET("/validate", func(c *gin.Context) {

		sub, _ := c.GetQuery("sub") // the user that wants to access a resource.
		obj, _ := c.GetQuery("obj") // the resource that is going to be accessed.
		act, _ := c.GetQuery("act") // the operation that the user performs on the resource.

		access := false

		if enforcer.Enforce(sub, obj, act) == true {
			access = true
		} else {
			access = false
		}

		c.JSON(200, gin.H{
			"message": access,
		})
	})

	r.GET("/save", func(c *gin.Context) {

		enforcer.AddPolicy("data2_admin", "data2", "delete")
		enforcer.SavePolicy()

		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.GET("/load", func(c *gin.Context) {

		enforcer.LoadPolicy()

		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.Run(":8080")
}
