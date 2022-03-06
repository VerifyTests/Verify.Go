package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Files struct {
}

// File an instance of the Files struct.
var File = Files{}

// GetFileNameWithoutExtension returns the name of the file without its extension.
func (f *Files) GetFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

// DeleteIfEmpty deletes the file at the path if it is empty
func (f *Files) DeleteIfEmpty(path string) {
	if f.Exists(path) && f.IsFileEmpty(path) {
		f.Delete(path)
	}
}

// IsFileEmpty checks if the file at the path is empty
func (f *Files) IsFileEmpty(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		panic(fmt.Sprintf("failed to get stat of the file at %s", path))
	}
	return info.Size() == 0
}

// Exists check if a file exists at the provided path
func (f *Files) Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// WriteText writes the provided text to the file
func (f *Files) WriteText(filePath string, text string) {
	var encodedText = []byte(text)
	err := ioutil.WriteFile(filePath, encodedText, 0600)
	if err != nil {
		panic(fmt.Sprintf("failed to write file at %s", filePath))
	}
}

// ReadText reads the content of the file as []byte
func (f *Files) ReadText(filePath string) []byte {
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

// Delete deletes the file at the path if exists
func (f *Files) Delete(path string) {
	if f.Exists(path) {
		err := os.Remove(path)
		if err != nil {
			panic(fmt.Sprintf("failed to delete the file at %s", path))
		}
	}
}

// Move moves the source file to the destination
func (f *Files) Move(sourcePath, destPath string) {
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

func (f *Files) close(file *os.File) {
	_ = file.Close()
}

// TryCreateFile creates a file at the path. If useEmptyStringForTextFiles is set and the
// file extension is text, creates an empty file with the provided name.
func (f *Files) TryCreateFile(path string, useEmptyStringForTextFiles bool) bool {
	Guard.AgainstEmpty(path)

	extension := f.GetFileExtension(path)
	if useEmptyStringForTextFiles && f.IsText(extension) {
		f.tryCreateDir(path)
		exists, err := f.FileOrDirectoryExists(path)
		if err != nil {
			panic(fmt.Sprintf("Failed to check if file exists: %s", err))
		}

		if exists {
			f.Delete(path)
		}

		f.WriteText(path, "")

		return true
	}

	return false
	//TODO: Implement default paths for other file types (jpeg)
}

func (f *Files) tryCreateDir(path string) {
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		panic(fmt.Sprintf("Failed to get base directory from %s", path))
	}

	if len(dir) > 0 {
		err := f.CreateDirectory(dir)
		if err != nil {
			panic(fmt.Sprintf("Failed to create directory: %s", dir))
		}
	}
}

// FileOrDirectoryExists checks if a file or directory exists
func (f *Files) FileOrDirectoryExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreateDirectory creates a directory
func (f *Files) CreateDirectory(directory string) error {
	if !f.Exists(directory) {
		return os.Mkdir(directory, 644)
	}
	return nil
}

// GetFileExtension returns the extension of the file
func (f *Files) GetFileExtension(extensionOrPath string) string {

	if !strings.ContainsAny(extensionOrPath, ".") {
		return extensionOrPath
	}

	ext := filepath.Ext(extensionOrPath)
	Guard.AgainstEmpty(ext)

	ext = strings.TrimPrefix(ext, ".")
	Guard.AgainstEmpty(ext)

	return ext
}

// IsText determines if an extension is a text file
func (f *Files) IsText(extensionOrPath string) bool {
	var extension = f.GetFileExtension(extensionOrPath)
	for _, txt := range textExtensions {
		if extension == txt {
			return true
		}
	}
	return false
}

// GetFileName returns the name of the file from the path
func (f *Files) GetFileName(path string) string {
	return filepath.Base(path)
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
