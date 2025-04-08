package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// 配置结构
type Config struct {
	AI                 string            `json:"AI"`
	Headers2           map[string]string `json:"headers2"`
	Suffixes           []string          `json:"suffixes"`
	AllowedRespHeaders []string          `json:"allowedRespHeaders"`
	APIKeys            struct {
		Kimi     string `json:"kimi"`
		DeepSeek string `json:"deepseek"`
		Qianwen  string `json:"qianwen"`
		HunYuan  string `json:"hunyuan"`
		Gpt      string `json:"gpt"`
		Glm      string `json:"glm"`
	} `json:"apiKeys"`
	RespBodyBWhiteList []string `json:"respBodyBWhiteList"`
}

// 全局配置变量
var conf Config

var Prompt = `{
    "role": "越权检测专家（专注HTTP响应语义分析）",
    "input_params": {
      "reqA": "原始请求对象（含URL/参数）",
      "responseA": "账号A正常请求的响应数据",
      "responseB": "替换为账号B凭证后的响应数据",
      "statusB": "账号B的HTTP状态码（优先级：403>500>200）",
      "dynamic_fields": ["timestamp", "nonce", "session_id", "uuid", "request_id"]
    },
    "analysis_flow": {
        "preprocessing": [
            "STEP1. 接口性质判断：通过reqA的URL/参数判断是否是/login /public等无需鉴权的接口",
            "STEP2. 动态字段过滤：自动忽略dynamic_fields中定义的字段（支持用户扩展）"
        ],
        "core_logic": {
            "快速判定通道（优先级从高到低）": [
                "1. 若resB.status_code为403/401 → 直接返回false",
                "2. 若resB包含'Access Denied'/'Unauthorized'等关键词 → 返回false",
                "3. 若resB为空(null/[]/{})且resA有数据 → 返回false",
                "4. 若resB包含resA的敏感字段（如user_id/email/balance） → 返回true",
                "5. 若resB.status_code为500 → 返回unknown"
            ],
            "深度分析模式（当快速通道未触发时执行）": {
                "结构对比": [
                    "a. 字段层级对比（使用JSON Path分析嵌套结构差异）",
                    "b. 关键字段匹配（如data/id/account相关字段的命名和位置）"
                ],
                "语义分析": [
                    "i. 数值型字段：检查是否符合同类型数据特征（如金额字段是否在合理范围）",
                    "ii. 文本型字段：检查命名规范是否一致（如用户ID是否为相同格式）"
                ]
            }
        }
    },
    "decision_tree": {
        "true_condition": [
            "非公共接口 && (结构相似度>80% || 包含敏感数据泄漏)",
            "关键业务字段（如订单号/用户ID）的命名和层级完全一致",
            "操作类接口返回success:true且结构相同（如修改密码成功）"
        ],
        "false_condition": [
            "公共接口（如验证码获取）",
            "结构差异显著（字段缺失率>30%）",
            "返回B账号自身数据（通过user_id、phone等字段判断）"
        ],
        "unknown_condition": [
            "结构部分匹配（50%-80%相似度）但无敏感数据",
            "返回数据为系统默认值（如false/null）",
            "存在加密/编码数据影响判断"
        ]
    },
    "output_spec": {
        "json": {
            "res": "\"true\", \"false\" 或 \"unknown\"",
            "reason": "简明技术结论（如'结构高度相似'/'包含敏感字段')，禁用模糊表述，例如：响应B包含用户A的邮箱信息 或者 结构相似度95%且无权限提示"
        }
    },
    "notes": [
        "仅输出 JSON 格式的结果，不添加任何额外文本或解释。",
        "确保 JSON 格式正确，便于后续处理。",
        "保持客观，仅根据响应内容进行分析。",
        "优先使用 HTTP 状态码、错误信息和数据结构匹配进行判断。",
        "支持用户提供额外的动态字段，提高匹配准确性。"
    ],
    "advanced_config": {
        "similarity_threshold": {
            "structure": 0.8,
            "content": 0.7
        },
        "sensitive_fields": [
            "password",
            "token",
            "phone",
            "id_card"
        ],
        "auto_retry": {
            "when": "检测到加密数据或非常规格式",
            "action": "建议提供解密方式后重新检测"
        }
    }
}
  `

// 加载配置文件
func loadConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &conf); err != nil {
		return err
	}

	return nil
}

// 获取配置
func GetConfig() Config {
	return conf
}

// 初始化配置
func init() {
	configPath := "./config.json" // 配置文件路径

	if err := loadConfig(configPath); err != nil {
		fmt.Printf("Error loading config file: %v\n", err)
		os.Exit(1)
	}
}
