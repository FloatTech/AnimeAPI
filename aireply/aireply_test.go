package aireply

import "testing"

func TestAireply(t *testing.T) {
	qyk := NewAIReply("青云客")
	t.Log("青云客测试:", qyk.Talk("你好", "椛椛"))
	xa := NewAIReply("小爱")
	t.Log("小爱测试:", xa.Talk("小米是垃圾", "椛椛"))
}
