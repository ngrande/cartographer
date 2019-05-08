package template

import (
	"os"
	"io/ioutil"
	"bufio"
	"log"
	"regexp"
	"strings"
	"path/filepath"
)

type template_entry struct {
	content string
	keywords []string
}

func ReadTemplateDir(dir string, begin_c rune, end_c rune) (map[string]template_entry, error) {
	
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed reading template dir '%s': %v", dir, err)
	}

	templates_map := make(map[string]template_entry)
	for _, f := range files {
		
		if f.IsDir() {
			sub, err := ReadTemplateDir(f.Name(), begin_c, end_c)
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
		templ, err := ReadTemplate(file_path, begin_c, end_c)
		if err != nil {
			log.Fatalf("Failed parsing template file '%s': %v", file_path, err)
		}

		// TODO key is only the name?!
		templates_map[f.Name()] = templ
	}

	return templates_map, nil
}

func ReadTemplate(file_path string, begin_c rune, end_c rune) (template_entry, error) {
	// read the template file and get all keywords

	regex_str := `\` + string(begin_c) + `.*` + `\` + string(end_c)
	log.Println("Regex str: ", regex_str)
	key_re, err := regexp.Compile(regex_str)
	if err != nil {
		log.Fatalf("Failed to compile regex for keyword analyzis (start: '%c' end: '%c'): %v", begin_c, end_c, err)
	}

	file, err := os.Open(file_path)
	if err != nil {
		log.Fatalf("Failed to read template file '%s': %v", file_path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var instance template_entry
	var buf strings.Builder
	for scanner.Scan() {
		txt := scanner.Text()
		buf.WriteString(txt)

		for _, str := range key_re.FindAllString(txt, -1) {
			log.Println("Found keyword:", str)

			instance.keywords = append(instance.keywords, []string{ str }...)
		}
	}

	instance.content = buf.String()

	return instance, nil
}
