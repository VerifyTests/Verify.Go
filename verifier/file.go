package verifier

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type fileHelper struct {
}

var file = fileHelper{}

func (f *fileHelper) getFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

func (f *fileHelper) deleteIfEmpty(path string) {
	if f.exists(path) && f.isFileEmpty(path) {
		f.delete(path)
	}
}

func (f *fileHelper) isFileEmpty(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		panic(fmt.Sprintf("failed to get stat of the file at %s", path))
	}
	return info.Size() == 0
}

func (f *fileHelper) exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func (f *fileHelper) writeText(filePath string, text string) {
	var encodedText = []byte(text)
	err := ioutil.WriteFile(filePath, encodedText, 0600)
	if err != nil {
		panic(fmt.Sprintf("failed to write file at %s", filePath))
	}
}

func (f *fileHelper) readText(filePath string) []byte {
	buf := bytes.NewBuffer(nil)
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to open file %s", filePath))
	}

	_, err = io.Copy(buf, file)
	if err != nil {
		panic(fmt.Sprintf("failed to read from file %s", filePath))
	}

	err = file.Close()
	if err != nil {
		panic(fmt.Sprintf("failed to close file %s", filePath))
	}

	return buf.Bytes()
}

func (f *fileHelper) getFileExtension(extensionOrPath string) string {

	if !strings.ContainsAny(extensionOrPath, ".") {
		return extensionOrPath
	}

	ext := filepath.Ext(extensionOrPath)
	guard.AgainstNullOrEmpty(ext)

	ext = strings.TrimPrefix(ext, ".")
	guard.AgainstNullOrEmpty(ext)

	return ext
}

func (f *fileHelper) isText(extensionOrPath string) bool {
	var extension = f.getFileExtension(extensionOrPath)
	for _, txt := range textExtensions {
		if extension == txt {
			return true
		}
	}
	return false
}

func (f *fileHelper) getFileName(path string) string {
	return filepath.Base(path)
}

func (f *fileHelper) delete(path string) {
	if f.exists(path) {
		err := os.Remove(path)
		if err != nil {
			panic(fmt.Sprintf("failed to delete the file at %s", path))
		}
	}
}

func (f *fileHelper) move(sourcePath, destPath string) {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		panic(fmt.Sprintf("couldn't open source file: %s", err))
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		f.close(inputFile)
		panic(fmt.Sprintf("couldn't open dest file: %s", err))
	}

	defer f.close(outputFile)
	_, err = io.Copy(outputFile, inputFile)
	f.close(inputFile)
	if err != nil {
		panic(fmt.Sprintf("writing to output file failed: %s", err))
	}

	err = os.Remove(sourcePath)
	if err != nil {
		panic(fmt.Sprintf("failed removing original file: %s", err))
	}
}

func (f *fileHelper) close(file *os.File) {
	_ = file.Close()
}

// From https://github.com/sindresorhus/text-extensions/blob/master/text-extensions.json
// contains list of text file extensions
//goland:noinspection SpellCheckingInspection
var textExtensions = []string{
	"ada",
	"adb",
	"ads",
	"applescript",
	"as",
	"asc",
	"ascii",
	"ascx",
	"asm",
	"asmx",
	"asp",
	"aspx",
	"atom",
	"au3",
	"awk",
	"bas",
	"bash",
	"bashrc",
	"bat",
	"bbcolors",
	"bcp",
	"bdsgroup",
	"bdsproj",
	"bib",
	"bowerrc",
	"c",
	"cbl",
	"cc",
	"cfc",
	"cfg",
	"cfm",
	"cfml",
	"cgi",
	"cjs",
	"clj",
	"cljs",
	"cls",
	"cmake",
	"cmd",
	"cnf",
	"cob",
	"code-snippets",
	"coffee",
	"coffeekup",
	"conf",
	"cp",
	"cpp",
	"cpt",
	"cpy",
	"crt",
	"cs",
	"csh",
	"cson",
	"csproj",
	"csr",
	"css",
	"csslintrc",
	"csv",
	"ctl",
	"curlrc",
	"cxx",
	"d",
	"dart",
	"dfm",
	"diff",
	"dof",
	"dpk",
	"dpr",
	"dproj",
	"dtd",
	"eco",
	"editorconfig",
	"ejs",
	"el",
	"elm",
	"emacs",
	"eml",
	"ent",
	"erb",
	"erl",
	"eslintignore",
	"eslintrc",
	"ex",
	"exs",
	"f",
	"f03",
	"f77",
	"f90",
	"f95",
	"fish",
	"for",
	"fpp",
	"frm",
	"fs",
	"fsproj",
	"fsx",
	"ftn",
	"gemrc",
	"gemspec",
	"gitattributes",
	"gitconfig",
	"gitignore",
	"gitkeep",
	"gitmodules",
	"go",
	"gpp",
	"gradle",
	"graphql",
	"groovy",
	"groupproj",
	"grunit",
	"gtmpl",
	"gvimrc",
	"h",
	"haml",
	"hbs",
	"hgignore",
	"hh",
	"hpp",
	"hrl",
	"hs",
	"hta",
	"htaccess",
	"htc",
	"htm",
	"html",
	"htpasswd",
	"hxx",
	"iced",
	"iml",
	"inc",
	"inf",
	"info",
	"ini",
	"ino",
	"int",
	"irbrc",
	"itcl",
	"itermcolors",
	"itk",
	"jade",
	"java",
	"jhtm",
	"jhtml",
	"js",
	"jscsrc",
	"jshintignore",
	"jshintrc",
	"json",
	"json5",
	"jsonld",
	"jsp",
	"jspx",
	"jsx",
	"ksh",
	"less",
	"lhs",
	"lisp",
	"log",
	"ls",
	"lsp",
	"lua",
	"m",
	"m4",
	"mak",
	"map",
	"markdown",
	"master",
	"md",
	"mdown",
	"mdwn",
	"mdx",
	"metadata",
	"mht",
	"mhtml",
	"mjs",
	"mk",
	"mkd",
	"mkdn",
	"mkdown",
	"ml",
	"mli",
	"mm",
	"mxml",
	"nfm",
	"nfo",
	"noon",
	"npmignore",
	"npmrc",
	"nuspec",
	"nvmrc",
	"ops",
	"pas",
	"pasm",
	"patch",
	"pbxproj",
	"pch",
	"pem",
	"pg",
	"php",
	"php3",
	"php4",
	"php5",
	"phpt",
	"phtml",
	"pir",
	"pl",
	"pm",
	"pmc",
	"pod",
	"pot",
	"prettierrc",
	"properties",
	"props",
	"pt",
	"pug",
	"purs",
	"py",
	"pyx",
	"r",
	"rake",
	"rb",
	"rbw",
	"rc",
	"rdoc",
	"rdoc_options",
	"resx",
	"rexx",
	"rhtml",
	"rjs",
	"rlib",
	"ron",
	"rs",
	"rss",
	"rst",
	"rtf",
	"rvmrc",
	"rxml",
	"s",
	"sass",
	"scala",
	"scm",
	"scss",
	"seestyle",
	"sh",
	"shtml",
	"sln",
	"sls",
	"spec",
	"sql",
	"sqlite",
	"sqlproj",
	"srt",
	"ss",
	"sss",
	"st",
	"strings",
	"sty",
	"styl",
	"stylus",
	"sub",
	"sublime-build",
	"sublime-commands",
	"sublime-completions",
	"sublime-keymap",
	"sublime-macro",
	"sublime-menu",
	"sublime-project",
	"sublime-settings",
	"sublime-workspace",
	"sv",
	"svc",
	"svg",
	"swift",
	"t",
	"tcl",
	"tcsh",
	"terminal",
	"tex",
	"text",
	"textile",
	"tg",
	"tk",
	"tmLanguage",
	"tmpl",
	"tmTheme",
	"tpl",
	"ts",
	"tsv",
	"tsx",
	"tt",
	"tt2",
	"ttml",
	"twig",
	"txt",
	"v",
	"vb",
	"vbproj",
	"vbs",
	"vcproj",
	"vcxproj",
	"vh",
	"vhd",
	"vhdl",
	"vim",
	"viminfo",
	"vimrc",
	"vm",
	"vue",
	"webapp",
	"webmanifest",
	"wsc",
	"x-php",
	"xaml",
	"xht",
	"xhtml",
	"xml",
	"xs",
	"xsd",
	"xsl",
	"xslt",
	"y",
	"yaml",
	"yml",
	"zsh",
	"zshrc",
}
