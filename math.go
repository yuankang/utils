package utils

import (
	"log"
	"math"
)

func MathTest() {
	i := 1
	log.Println(math.Abs(float64(i))) //取到绝对值
	log.Println(math.Ceil(3.8))       //向上取整
	log.Println(math.Floor(3.6))      //向下取整
	log.Println(math.Mod(11, 3))      //取余数 11%3 效果一样
	log.Println(math.Modf(3.22))      //取整数跟小数
	log.Println(math.Pow(3, 2))       //X 的 Y次方  乘方
	log.Println(math.Pow10(3))        //10的N次方 乘方
	log.Println(math.Sqrt(9))         //开平方  3
	log.Println(math.Cbrt(8))         //开立方  2
	log.Println(math.Pi)              //π
	log.Println(math.Round(4.2))      //四舍五入

	log.Println(math.IsNaN(3.4))      //false   报告f是否表示一个NaN（Not A Number）值。
	log.Println(math.Trunc(1.999999)) //1    返回整数部分（的浮点值）。
	log.Println(math.Max(-1.3, 0))    //0   返回x和y中最大值
	log.Println(math.Min(-1.3, 0))    //-1.3  返回x和y中最小值
	log.Println(math.Dim(-12, -19))   //7 函数返回x-y和0中的最大值
	log.Println(math.Dim(-12, 19))    //0 函数返回x-y和0中的最大值
	log.Println(math.Cbrt(8))         //2  返回x的三次方根
	log.Println(math.Hypot(3, 4))     //5  返回Sqrt(p*p + q*q)
	log.Println(math.Pow(2, 8))       //256  返回x^y
}

/*
1
4
3
2
3 0.2200000000000002
9
1000
3
2
3.141592653589793
4
false
1
0
-1.3
7
0
2
5
256
*/
