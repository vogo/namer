package data

// SanCaiJiXiong 三才配置吉凶表，key 为 "天五行人五行地五行"
var SanCaiJiXiong = map[string]JiXiong{
	// 天格木
	"木木木": DaJi, "木木火": DaJi, "木木土": DaJi, "木木金": XiongDuo, "木木水": JiDuo,
	"木火木": DaJi, "木火火": ZhongJi, "木火土": DaJi, "木火金": XiongDuo, "木火水": DaXiong,
	"木土木": DaXiong, "木土火": ZhongJi, "木土土": Ji, "木土金": JiDuo, "木土水": DaXiong,
	"木金木": DaXiong, "木金火": DaXiong, "木金土": XiongDuo, "木金金": DaXiong, "木金水": DaXiong,
	"木水木": DaJi, "木水火": XiongDuo, "木水土": XiongDuo, "木水金": DaJi, "木水水": DaJi,
	// 天格火
	"火木木": DaJi, "火木火": DaJi, "火木土": DaJi, "火木金": XiongDuo, "火木水": ZhongJi,
	"火火木": DaJi, "火火火": ZhongJi, "火火土": DaJi, "火火金": DaXiong, "火火水": DaXiong,
	"火土木": JiDuo, "火土火": DaJi, "火土土": DaJi, "火土金": DaJi, "火土水": JiDuo,
	"火金木": DaXiong, "火金火": DaXiong, "火金土": JiXiongBan, "火金金": DaXiong, "火金水": DaXiong,
	"火水木": XiongDuo, "火水火": DaXiong, "火水土": DaXiong, "火水金": DaXiong, "火水水": DaXiong,
	// 天格土
	"土木木": ZhongJi, "土木火": ZhongJi, "土木土": XiongDuo, "土木金": DaXiong, "土木水": XiongDuo,
	"土火木": DaJi, "土火火": DaJi, "土火土": DaJi, "土火金": JiDuo, "土火水": DaXiong,
	"土土木": ZhongJi, "土土火": DaJi, "土土土": DaJi, "土土金": DaJi, "土土水": XiongDuo,
	"土金木": XiongDuo, "土金火": XiongDuo, "土金土": DaJi, "土金金": DaJi, "土金水": DaJi,
	"土水木": XiongDuo, "土水火": DaXiong, "土水土": DaXiong, "土水金": JiXiongBan, "土水水": DaXiong,
	// 天格金
	"金木木": XiongDuo, "金木火": XiongDuo, "金木土": XiongDuo, "金木金": DaXiong, "金木水": XiongDuo,
	"金火木": XiongDuo, "金火火": JiXiongBan, "金火土": JiXiongBan, "金火金": DaXiong, "金火水": DaXiong,
	"金土木": ZhongJi, "金土火": DaJi, "金土土": DaJi, "金土金": DaJi, "金土水": JiDuo,
	"金金木": DaXiong, "金金火": DaXiong, "金金土": DaJi, "金金金": ZhongJi, "金金水": ZhongJi,
	"金水木": DaJi, "金水火": XiongDuo, "金水土": Ji, "金水金": DaJi, "金水水": ZhongJi,
	// 天格水
	"水木木": DaJi, "水木火": DaJi, "水木土": DaJi, "水木金": XiongDuo, "水木水": DaJi,
	"水火木": ZhongJi, "水火火": DaXiong, "水火土": XiongDuo, "水火金": DaXiong, "水火水": DaXiong,
	"水土木": DaXiong, "水土火": ZhongJi, "水土土": ZhongJi, "水土金": ZhongJi, "水土水": DaXiong,
	"水金木": XiongDuo, "水金火": XiongDuo, "水金土": DaJi, "水金金": ZhongJi, "水金水": DaJi,
	"水水木": DaJi, "水水火": DaXiong, "水水土": DaXiong, "水水金": DaJi, "水水水": ZhongJi,
}
