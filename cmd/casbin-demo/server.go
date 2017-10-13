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
		"user=tlantic password=....! host=... sslmode=disable dbname=authz",
		true)

	enforcer = casbin.NewEnforcer(m, adapter)

	enforcer.LoadPolicy()

	for i:=0 ; i<1000; i++{
		enforcer.AddPolicy( fmt.Sprintf("alice-%d",i), "data1", "read")
	}

	/*enforcer.AddPolicy( "alice", "data1", "read")
	enforcer.AddPolicy("data2_admin", "data2", "read")
	enforcer.AddPolicy("data2_admin", "data2", "write")
	//enforcer.SavePolicy()
	enforcer.AddGroupingPolicy("alice", "data2_admin")
	enforcer.SavePolicy()

	fmt.Printf("get policy: %v\n", enforcer.GetPolicy()) */
}

func main() {
	LoadCasbin()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		sub := "alice" // the user that wants to access a resource.
		obj := "data2" // the resource that is going to be accessed.
		act := "read" // the operation that the user performs on the resource.

		if enforcer.Enforce(sub, obj, act) == true {
			fmt.Println("aa")
		} else {
			fmt.Println("bb")
		}

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ping2", func(c *gin.Context) {

		enforcer.AddPolicy("data2_admin", "data2", "read")
		enforcer.SavePolicy()

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
}
