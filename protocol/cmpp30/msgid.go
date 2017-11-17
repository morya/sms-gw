package cmpp30

// msgid 算法

/*
MSGID 信息标识

生成算法如下：
	采用64位（8字节）的整数：
		（1）时间（格式为 MMDDHHMMSS，即月日时分秒）：
			bit64~bit39，
			其中
			bit64~bit61：月份的二进制表示；
			bit60~bit56：日的二进制表示；
			bit55~bit51：小时的二进制表示；
			bit50~bit45：分的二进制表示；
			bit44~bit39：秒的二进制表示；
		（2）短信网关代码：
			bit38~bit17， 把短信网关的代码转换为整数填写到该字段中。
		（3）序列号：
			bit16~bit1，顺序增加， 步长为1，循环使用。

	各部分如不能填满，左补零，右对齐。
*/
