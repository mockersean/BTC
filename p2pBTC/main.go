package main

import "fmt"

func main() {
	//bc := NewBlockChain()
	//cli := CLI{bc}
	cli := CLI{}
	cli.Run()


	
	//m := pass()
	//for k ,v := range m{
	//	fmt.Printf("key = %s,value =%v\n",k,v)
	//}

}

type Student struct {
	Name string
	Age int
}

func pass() map[string]*Student  {
	m := make(map[string]*Student)
	stu := []Student{{"sean",12},{"hyq",10}}
	for _,v := range stu{
		fmt.Printf("地址%p\n",&v)
		m[v.Name] = &v
	}
	return m
}