package srcscan

import (
	"github.com/kr/pretty"
	"go/build"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	Default.PathIndependent = true
	Default.Base = "testdata"
	type scanTest struct {
		config *Config
		dir    string
		units  []Unit
	}
	tests := []scanTest{
		{
			dir: "testdata",
			units: []Unit{
				&BowerComponent{
					Dir:       "bower",
					BowerJSON: []byte(`{"name":"foo","dependencies":{"baz":"1.0.0"}}`),
				},
				&GoPackage{
					Package: build.Package{
						Dir:            "go",
						Name:           "mypkg",
						ImportPath:     "github.com/sourcegraph/srcscan/testdata/go",
						GoFiles:        []string{"a.go", "b.go"},
						Imports:        []string{},
						ImportPos:      nil,
						TestGoFiles:    []string{"a_test.go"},
						TestImports:    []string{},
						TestImportPos:  nil,
						XTestGoFiles:   []string{"b_test.go"},
						XTestImports:   []string{},
						XTestImportPos: nil,
					},
				},
				&GoPackage{
					Package: build.Package{
						Dir:            "go/cmd/mycmd",
						Name:           "main",
						ImportPath:     "github.com/sourcegraph/srcscan/testdata/go/cmd/mycmd",
						GoFiles:        []string{"mycmd.go"},
						Imports:        []string{},
						ImportPos:      nil,
						TestGoFiles:    nil,
						TestImports:    []string{},
						TestImportPos:  nil,
						XTestGoFiles:   nil,
						XTestImports:   []string{},
						XTestImportPos: nil,
					},
				},
				&GoPackage{
					Package: build.Package{
						Dir:            "go/qux",
						Name:           "qux",
						ImportPath:     "github.com/sourcegraph/srcscan/testdata/go/qux",
						GoFiles:        []string{"qux.go"},
						Imports:        []string{},
						ImportPos:      nil,
						TestGoFiles:    nil,
						TestImports:    []string{},
						TestImportPos:  nil,
						XTestGoFiles:   nil,
						XTestImports:   []string{},
						XTestImportPos: nil,
					},
				},
				&JavaProject{
					Dir:              "java-maven",
					ProjectClasspath: "target/classes",
					SrcFiles:         []string{"src/main/java/foo/Foo.java"},
					TestFiles:        []string{"src/test/java/bar/Bar.java"},
				},
				&NPMPackage{
					Dir:            "npm",
					PackageJSON:    []byte(`{"name":"mypkg"}`),
					LibFiles:       []string{"a.js", "lib/a.js"},
					TestFiles:      []string{"a_test.js", "test/b.js", "test/c_test.js"},
					VendorFiles:    []string{"example/bower_components/foo/foo.js", "vendor/a.js"},
					GeneratedFiles: []string{"a.min.js", "dist/a.js"},
				},
				&NPMPackage{
					Dir:         "npm/subpkg",
					PackageJSON: []byte(`{"name":"subpkg"}`),
					LibFiles:    []string{"a.js"},
				},
				&PythonModule{"python/myscript.py"},
				&PythonPackage{"python/mypkg"},
				&RubyApp{
					Dir:       "ruby/sample_app",
					SrcFiles:  []string{"app/app.rb"},
					TestFiles: nil,
				},
				&RubyGem{
					Dir:         "ruby/sample_gem",
					Name:        "sample_ruby_gem",
					GemSpecFile: "sample_ruby_gem.gemspec",
					SrcFiles:    []string{"lib/sample_ruby_gem.rb"},
					TestFiles:   []string{"spec/my_spec.rb", "test/qux.rb", "test/test_foo.rb"},
				},
			},
		},
		{
			config: &Config{
				PathIndependent: true,
				Base:            "testdata",
				Profiles: []Profile{
					{
						Name:         "Python package",
						TopLevelOnly: false,
						Dir:          FileInDir{"__init__.py"},
						File:         FileHasSuffix{".py"},
						Unit: func(abspath, relpath string, config Config, info os.FileInfo) Unit {
							if info.IsDir() {
								return &PythonPackage{relpath}
							} else {
								return &PythonModule{relpath}
							}
						},
					},
				},
			},
			dir: "testdata/python",
			units: []Unit{
				&PythonModule{"python/mypkg/__init__.py"},
				&PythonModule{"python/mypkg/a.py"},
				&PythonModule{"python/mypkg/qux/__init__.py"},
				&PythonModule{"python/myscript.py"},
				&PythonPackage{"python/mypkg"},
				&PythonPackage{"python/mypkg/qux"},
			},
		},
	}
	for _, test := range tests {
		// Use default config if config is nil.
		var config Config
		if test.config != nil {
			config = *test.config
		} else {
			config = Default
		}

		units, err := config.Scan(test.dir)
		if err != nil {
			t.Errorf("got error %q", err)
			continue
		}

		sort.Sort(Units(units))
		sort.Sort(Units(test.units))
		if !reflect.DeepEqual(test.units, units) {
			t.Errorf("units:\n%v", pretty.Diff(test.units, units))
			if len(test.units) == len(units) {
				for i := range test.units {
					if !reflect.DeepEqual(test.units[i], units[i]) {
						t.Errorf("units[%d]:\n%v", i, strings.Join(pretty.Diff(test.units[i], units[i]), "\n"))
					}
				}
			}
		}
	}
}
