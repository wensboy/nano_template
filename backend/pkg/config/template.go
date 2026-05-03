package config

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/adrg/frontmatter"

	"example.com/nano_template/pkg/util"
)

type (
	TemplateFrontMatter struct {
		Role string `yaml:"role" toml:"role" json:"role"`
	}
	Template struct {
		FrontMatter TemplateFrontMatter
		Content     []byte `json:"content"`
	}
	TemplateMeta struct {
		Id     string   `yaml:"id"`
		Tfiles []string `yaml:"tFiles"`
	}
	TemplateConfig struct {
		BasePath string         `yaml:"basePath"`
		Suffix   string         `yaml:"suffix"`
		Tlist    []TemplateMeta `yaml:"tList"`
	}
)

func DefaultTemplateConfig() TemplateConfig {
	return TemplateConfig{
		BasePath: "",
		Suffix:   ".md",
		Tlist:    []TemplateMeta{},
	}
}

func MapTemplates(cfg *TemplateConfig) {
	if cfg == nil {
		util.Warn("Template config is nil, skipping template loading")
		return
	}

	basePath := filepath.Clean(cfg.BasePath)
	if !util.DirExists(basePath) {
		util.Warn(fmt.Sprintf("Template base path %s does not exist, skipping template loading", basePath))
		return
	}

	jobs := buildTemplateLoadJobs(cfg, basePath)
	if len(jobs) == 0 {
		getGlobalTemplatePool().Reset(map[string]Template{})
		return
	}

	go func(tasks []templateLoadJob) {
		pool := util.NewPool(util.WithMaxWorkers(len(tasks)))
		defer pool.Close()

		results := util.Map(
			pool,
			context.Background(),
			tasks,
			func(_ context.Context, task templateLoadJob) (Template, error) {
				return readTemplate(task.Path)
			},
		)

		templates := make(map[string]Template, len(results))
		failed := 0
		for _, result := range results {
			if result.Err != nil {
				failed++
				util.Warn(fmt.Sprintf("Load template %s from %s failed: %v", result.Input.Key, result.Input.Path, result.Err))
				continue
			}

			templates[result.Input.Key] = result.Value
		}

		getGlobalTemplatePool().Reset(templates)
		util.Info(fmt.Sprintf("Template mapping finished: loaded=%d failed=%d", len(templates), failed))
	}(jobs)
}

type templateLoadJob struct {
	Key  string
	Path string
}

func buildTemplateLoadJobs(cfg *TemplateConfig, basePath string) []templateLoadJob {
	suffix := cfg.Suffix
	if suffix == "" {
		suffix = DefaultTemplateConfig().Suffix
	}

	jobs := make([]templateLoadJob, 0)
	for _, meta := range cfg.Tlist {
		if meta.Id == "" {
			continue
		}
		for _, tFile := range meta.Tfiles {
			if tFile == "" {
				continue
			}

			jobs = append(jobs, templateLoadJob{
				Key:  fmt.Sprintf("%s_%s", meta.Id, tFile),
				Path: filepath.Join(basePath, tFile+suffix),
			})
		}
	}

	return jobs
}

func readTemplate(path string) (Template, error) {
	data, err := util.ReadBytes(path)
	if err != nil {
		return Template{}, err
	}

	var meta TemplateFrontMatter
	content, err := frontmatter.Parse(bytes.NewReader(data), &meta)
	if err != nil {
		return Template{}, err
	}

	return Template{
		FrontMatter: meta,
		Content:     append([]byte(nil), content...),
	}, nil
}

type (
	templatePool struct {
		mu        sync.RWMutex
		templates map[string]Template
	}
)

var (
	globalTemplatePool     *templatePool
	globalTemplatePoolOnce sync.Once
)

func newTemplatePool() *templatePool {
	return &templatePool{
		templates: make(map[string]Template),
	}
}

func getGlobalTemplatePool() *templatePool {
	globalTemplatePoolOnce.Do(func() {
		globalTemplatePool = newTemplatePool()
	})
	return globalTemplatePool
}

func cloneTemplate(t Template) Template {
	cloned := t
	if t.Content != nil {
		cloned.Content = append([]byte(nil), t.Content...)
	}
	return cloned
}

func (p *templatePool) Set(id string, tmpl Template) {
	if id == "" {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.templates[id] = cloneTemplate(tmpl)
}

func (p *templatePool) Get(id string) (Template, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	tmpl, ok := p.templates[id]
	if !ok {
		return Template{}, false
	}

	return cloneTemplate(tmpl), true
}

func (p *templatePool) Delete(id string) {
	if id == "" {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.templates, id)
}

func (p *templatePool) Reset(templates map[string]Template) {
	next := make(map[string]Template, len(templates))
	for id, tmpl := range templates {
		if id == "" {
			continue
		}
		next[id] = cloneTemplate(tmpl)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.templates = next
}

func (p *templatePool) List() map[string]Template {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]Template, len(p.templates))
	for id, tmpl := range p.templates {
		result[id] = cloneTemplate(tmpl)
	}
	return result
}

// SetTemplate stores a template in the package-global pool.
func SetTemplate(id string, tmpl Template) {
	getGlobalTemplatePool().Set(id, tmpl)
}

// GetTemplate returns a template from the package-global pool by id.
func GetTemplate(id string) (Template, bool) {
	return getGlobalTemplatePool().Get(id)
}

// DeleteTemplate removes a template from the package-global pool by id.
func DeleteTemplate(id string) {
	getGlobalTemplatePool().Delete(id)
}

// ResetTemplatePool replaces the full package-global pool state.
func ResetTemplatePool(templates map[string]Template) {
	getGlobalTemplatePool().Reset(templates)
}

// ListTemplates returns a snapshot of all templates in the package-global pool.
func ListTemplates() map[string]Template {
	return getGlobalTemplatePool().List()
}
