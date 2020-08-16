package teaconfigs

// RewriteList接口
type RewriteListInterface interface {
	// 校验
	ValidateRewriteRules() error

	// 取得所有的Rewrite
	AllRewriteRules() []*RewriteRule

	// 根据ID查找Rewrite
	FindRewriteRule(rewriteId string) *RewriteRule

	// 添加Rewrite
	AddRewriteRule(rewrite *RewriteRule)

	// 删除Rewrite
	RemoveRewriteRule(rewriteId string)
}

// RewriteList定义
type RewriteList struct {
	Rewrite []*RewriteRule `yaml:"rewrite" json:"rewrite"`
}

// 校验
func (this *RewriteList) ValidateRewriteRules() error {
	for _, r := range this.Rewrite {
		err := r.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// 取得所有的Rewrite
func (this *RewriteList) AllRewriteRules() []*RewriteRule {
	if this.Rewrite == nil {
		return []*RewriteRule{}
	}
	return this.Rewrite
}

// 根据ID查找Rewrite
func (this *RewriteList) FindRewriteRule(rewriteId string) *RewriteRule {
	for _, r := range this.Rewrite {
		if r.Id == rewriteId {
			r.Validate()
			return r
		}
	}
	return nil
}

// 添加Rewrite
func (this *RewriteList) AddRewriteRule(rewrite *RewriteRule) {
	this.Rewrite = append(this.Rewrite, rewrite)
}

// 删除Rewrite
func (this *RewriteList) RemoveRewriteRule(rewriteId string) {
	result := []*RewriteRule{}
	for _, r := range this.Rewrite {
		if r.Id == rewriteId {
			continue
		}
		result = append(result, r)
	}
	this.Rewrite = result
}
