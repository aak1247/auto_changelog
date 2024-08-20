package gitope

import (
	"fmt"
	"github.com/aak1247/gchangelog/configs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ChangeLog struct {
	Version string
	Head    *object.Commit
	Groups  map[string][]*object.Commit
}

func (c *ChangeLog) ParseCommits(commits []*object.Commit) {
	// 分组
	for _, commit := range commits {
		t, _ := ParseCommitMessage(commit)
		if g, ok := c.Groups[t]; ok == false {
			g = append(g, commit)
			c.Groups[t] = g
		} else {
			c.Groups[t] = []*object.Commit{commit}
		}
	}
}

func (c *ChangeLog) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("## %s    <sub>[%s](%s) - [%s](%s) [CI](%s)</sub>\n\n", c.Version,
		c.Head.Author.When.Format("2006-01-02"), GetTagUrl(configs.BaseUrl, configs.Project, c.Version),
		c.Head.Hash.String()[:8], GetCommitUrl(configs.BaseUrl, configs.Project, c.Version),
		GetTagPipelineUrl(configs.BaseUrl, configs.Project, c.Version)))
	for _, k := range configs.Types {
		if v, ok := c.Groups[k]; ok {
			s.WriteString(fmt.Sprintf("### %s\n", k))
			for _, v := range v {
				// TODO
				s.WriteString(fmt.Sprintf("- [%s ( %s by %s )](%s)\n", v.Message, v.Hash.String()[:8], v.Author.Name, GetCommitUrl(configs.BaseUrl, configs.Project, v.Hash.String())))
			}
		}
	}
	return s.String()
}

func FindCommits(tag2 *plumbing.Reference, tag1 *plumbing.Reference, r *git.Repository) []*object.Commit {
	var tag1Hash, tag2Hash string
	var options = &git.LogOptions{
		From: tag1.Hash(),
	}
	tag1Hash = tag1.Hash().String()
	var startTime, endTime time.Time
	var tag1Head, tag2Head *object.Tag
	var commitHead1, commitHead2 *object.Commit
	tag1Head, err := r.TagObject(tag1.Hash())
	if err != nil {
		if err == plumbing.ErrObjectNotFound {
			commitHead1, err = r.CommitObject(tag1.Hash())
			if err != nil {
				log.Println("failed to find tag ref ", err)
				endTime = time.Now()
			}
		}
		endTime = commitHead1.Author.When
	} else {
		endTime = tag1Head.Tagger.When
	}
	options.Until = &(endTime)
	if tag2 != nil {
		// 有旧tag
		tag2Hash = tag2.Hash().String()
		tag2Head, err = r.TagObject(tag2.Hash())
		if err != nil {
			if err == plumbing.ErrObjectNotFound {
				commitHead2, err = r.CommitObject(tag2.Hash())
				if err != nil {
					log.Println("failed to find tag ref ", err)
					startTime = time.UnixMilli(0)
				}
			}
			startTime = commitHead2.Author.When
		} else {
			startTime = tag2Head.Tagger.When
		}
		options.Since = &(startTime)
	}
	// 遍历两个tag中间的log, 通过hash
	logIter, err := r.Log(options)
	if err != nil {
		panic(err)
	}
	commits := make([]*object.Commit, 0)
	var start, end bool
	for {
		commit, err := logIter.Next()
		if err != nil {
			break
		}
		if end {
			break
		}
		if commit.Hash.String() == tag1Hash {
			// 开始
			start = true
		}
		if commit.Hash.String() == tag2Hash {
			// 结束
			end = true
		}
		if start && !end {
			commits = append(commits, commit)
		}
		if end {
			break
		}
	}
	return commits
}

func FindTag(err error, r *git.Repository) (*plumbing.Reference, *plumbing.Reference, error) {
	// 先拿到最近的两个tag
	var tag1, tag2 *plumbing.Reference
	var tag1Name string
	tagIter, err := r.Tags()
	if err != nil {
		panic(err)
	}
	tag1, err = tagIter.Next()
	tag1Name = TagName(tag1)
	// 遍历找到最后两个
	for {
		tagN, err := tagIter.Next()
		if err != nil || tagN == nil {
			break
		}
		tagNName := TagName(tagN)
		if VersionCompare(tagNName, tag1Name) > 0 {
			tag1Name = TagName(tag1)
			tag2 = tag1
			tag1 = tagN
		}
	}

	if tag1 == nil && tag2 == nil {
		// 没有tag
		panic("no tag found")
	}
	if err != nil {
		// 报错
		panic(err)
	}
	return tag1, tag2, err
}

func ParseCommitMessage(commit *object.Commit) (typ string, msg string) {
	fullMsg := commit.Message
	s, _ := commit.Stats()
	s.String()
	// 根据configs.Types 解析类型和实际配置
	for _, t := range configs.Types {
		upper := strings.ToUpper(t)
		camel := strings.ToTitle(t)
		if strings.HasPrefix(fullMsg, camel) || strings.HasPrefix(fullMsg, t) || strings.HasPrefix(fullMsg, upper) {
			return t, MakeCommitMessage(commit)
		}
	}
	return "other", MakeCommitMessage(commit)
}

func MakeCommitMessage(commit *object.Commit) string {
	msg := fmt.Sprintf("%s (by %s)", commit.Message, commit.Author)
	return msg
}

func TagName(ref *plumbing.Reference) string {
	return strings.TrimPrefix(ref.Name().String(), "refs/tags/")
}

func GetProjectPath(r *git.Repository) string {
	remotes, err := r.Remotes()
	if err != nil {
		panic(err)
	}
	var fullUrl string
	//var remoteUrl string
	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			fullUrl = remote.String()
		}
	}
	s := strings.Split(fullUrl, "//")[1]
	repo := strings.SplitN(s, "/", 2)[1]
	path := strings.Split(repo, ".")[0]
	return path
}

func GetBaseUrl(r *git.Repository) string {
	remotes, err := r.Remotes()
	if err != nil {
		panic(err)
	}
	var fullUrl string
	//var remoteUrl string
	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			fullUrl = remote.String()
		}
	}
	s := strings.Split(fullUrl, "//")[1]
	baseUrl := strings.Split(s, "/")[0]
	if strings.Contains(baseUrl, "@") {
		// 去掉@
		baseUrl = strings.Split(baseUrl, "@")[1]
	}
	if strings.Contains(baseUrl, ":") {
		// 去掉:
		baseUrl = strings.Split(baseUrl, ":")[0]
	}
	return baseUrl
}

func GetCommitUrl(base, project, hash string) string {
	return fmt.Sprintf("https://%s/%s/-/commit/%s", base, project, hash)
}

func GetTagUrl(base, project, tagName string) string {
	return fmt.Sprintf("https://%s/%s/-/tags/%s", base, project, tagName)
}

func GetTagPipelineUrl(base, project, tagName string) string {
	return fmt.Sprintf("https://%s/%s/pipelines?page=1&scope=tags&ref=%s", base, project, tagName)
}

func VersionCompare(v1, v2 string) int {
	makeup := func(s []string) []string {
		if len(s) < 3 {
			for i := len(s) - 1; i < 3; i++ {
				s = append(s, "0")
			}
		}
		return s
	}
	rep := regexp.MustCompile(`[-_+]`)
	// 根据语义化版本号比较两个版本的大小
	v1 = strings.TrimPrefix(v1, "v")
	v1 = strings.TrimPrefix(v1, "V")
	v2 = strings.TrimPrefix(v2, "v")
	v2 = strings.TrimPrefix(v2, "V")
	// 按major minor patch 分割，然后分别比较
	s1 := makeup(strings.Split(v1, "."))
	major1, minor1, patch1 := s1[0], s1[1], s1[2]
	s2 := makeup(strings.Split(v2, "."))
	major2, minor2, patch2 := s2[0], s2[1], s2[2]
	if major1 != major2 {
		num1, _ := strconv.Atoi(major1)
		num2, _ := strconv.Atoi(major2)
		return num1 - num2
	}
	if minor1 != minor2 {
		num1, _ := strconv.Atoi(minor1)
		num2, _ := strconv.Atoi(minor2)
		return num1 - num2
	}
	if patch1 != patch2 {
		// 去掉后缀
		s1 := rep.Split(patch1, 2)
		s2 := rep.Split(patch2, 2)
		if s1[0] == s2[0] {
			if len(s1) == 2 && len(s1[1]) == 2 {
				return strings.Compare(s1[1], s2[1])
			} else {
				return len(s1) - len(s2)
			}
		} else {
			num1, _ := strconv.Atoi(s1[0])
			num2, _ := strconv.Atoi(s2[0])
			return num1 - num2
		}
	}
	return 0
}
