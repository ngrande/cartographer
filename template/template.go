package template

import (
	"os"
	"io/ioutil"
	"bufio"
	"log"
	"regexp"
	"strings"
	"fmt"
	"path/filepath"
)

type TemplateEntry struct {
	content string
	keywords []string
	key_sep rune
}

// Template naming scheme is important
// it will affect the way files are parsed and saved
// Template name excample: default.html
// File using the template will have to prepend the whole template name like:
// index.default.html => this will use default.html as template and later
// this file will be saved / read as index.html


func ReplaceWithTemplate(name string, content string, templates_map map[string]TemplateEntry) (string, string, bool) {

	for k, v := range templates_map {
		if !strings.HasSuffix(name, k) {
			continue
		}

		prefix := name[:len(strings.TrimSuffix(name, k)) - 1]
		name_res := prefix + filepath.Ext(k)
		content_res := v.content

		for _, keyword := range v.keywords {
			// quirk here:
			// set the "s" flag to match multiline!
			regex_str := fmt.Sprintf(`(?s)<\%c%s\%c>(?P<content>.*)<\%c%s\%c>`, v.key_sep, keyword, v.key_sep, v.key_sep, keyword, v.key_sep)

//			log.Println(regex_str)
			replace_re := regexp.MustCompile(regex_str)

			re_res := replace_re.FindStringSubmatch(content)
//			log.Println(re_res)
			if len(re_res) <= 1 {
				log.Fatalf("Not every keyword of the template was implemented! Missing '%s' in file '%s'", keyword, name)
				continue
			}
			// match - replace!
			replace_str := fmt.Sprintf("%c%s%c", v.key_sep, keyword, v.key_sep)
//			log.Printf("Replacing '%s' with '%s'", replace_str, re_res[1])
			content_res = strings.ReplaceAll(content_res, replace_str, re_res[1])
		}

		// done
		return name_res, content_res, true
	}

	return name, content, false
}

func ReadTemplateDir(dir string, sep rune) (map[string]TemplateEntry, error) {
	
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed reading template dir '%s': %v", dir, err)
	}

	templates_map := make(map[string]TemplateEntry)
	for _, f := range files {
		
		if f.IsDir() {
			sub, err := ReadTemplateDir(f.Name(), sep)
			if err != nil {
				log.Fatalf("Failed rading subdir '%s': %v", f.Name(), err)
			}
			for k, v := range sub {
				// TODO what about duplicates?!
				templates_map[k] = v
			}
			continue
		}

		file_path := filepath.Join(dir, f.Name())
		templ, err := ReadTemplate(file_path, sep)
		if err != nil {
			log.Fatalf("Failed parsing template file '%s': %v", file_path, err)
		}

		// TODO key is only the name?!
		templates_map[f.Name()] = templ
	}

	return templates_map, nil
}

func ReadTemplate(file_path string, sep rune) (TemplateEntry, error) {
	// read the template file and get all keywords

	regex_str := `\` + string(sep) + `.*` + `\` + string(sep)
	log.Println("Regex str: ", regex_str)
	key_re := regexp.MustCompile(regex_str)

	file, err := os.Open(file_path)
	if err != nil {
		log.Fatalf("Failed to read template file '%s': %v", file_path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var instance TemplateEntry
	var buf strings.Builder
	for scanner.Scan() {
		txt := scanner.Text()
		buf.WriteString(txt)

		for _, str := range key_re.FindAllString(txt, -1) {
			log.Println("Found keyword:", str)

			instance.keywords = append(instance.keywords, []string{ str[1:len(str) - 1] }...)
		}
	}

	instance.content = buf.String()
	instance.key_sep = sep

	return instance, nil
}
