# mystake
 股票涨跌幅监控

## 使用方法:
打包后命令行使用:
mystake monitor up [-i] -f <文件名> -u <涨幅提醒值> -d <跌幅提醒值>
示例: mystake monitor up -i true -f /Users/Chen/Desktop/Table.txt -u 3.0 -d -5.0

-i: 默认监控时间: 工作日 9:40-11:30, 13:00-14:55, -i 参数将忽略时间, 可用于调试. 
文件: 文本文件, 内容为要监控的股票代码, 每行一个, 默认获取前 6 位数字.示例:
SZ000936 华西股份
SH603767 中马传动
程序会自动识别000936,603767.
