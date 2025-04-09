package goflint

func NewSparkSubmitCmd(base *SparkApp) *SparkSubmitCmd {
	// если base == nil, начинаем с пустого
	cmd := &SparkSubmitCmd{
		appPath: "",
		opts:    make(map[string]string),
		flags:   make(map[string]bool),
		args:    []string{},
	}
	if base != nil {
		cmd.appPath = base.appPath
		for k, v := range base.opts {
			cmd.opts[k] = v
		}
		for k, v := range base.flags {
			cmd.flags[k] = v
		}
		cmd.args = slices.Clone(base.args)
	}
	return cmd
}

type SparkSubmitCmd struct {
	appPath string
	args    []string
	opts    map[string]string
	flags   map[string]bool
}

func (s *SparkSubmitCmd) WithApp(appPath string) *SparkSubmitCmd {
	s.appPath = appPath
	return s
}

func (s *SparkSubmitCmd) WithArg(arg string) *SparkSubmitCmd {
	s.args = append(s.args, arg)
	return s
}

func (s *SparkSubmitCmd) WithOpt(key, value string) *SparkSubmitCmd {
	s.opts[key] = value
	return s
}

func (s *SparkSubmitCmd) WithFlag(flag string) *SparkSubmitCmd {
	s.flags[flag] = true
	return s
}

// Build returns a new immutable SparkApp
func (s *SparkSubmitCmd) Build() *SparkApp {
	newOpts := make(map[string]string)
	for k, v := range s.opts {
		newOpts[k] = v
	}
	newFlags := make(map[string]bool)
	for k, v := range s.flags {
		newFlags[k] = v
	}
	return &SparkApp{
		appPath: s.appPath,
		args:    slices.Clone(s.args),
		opts:    newOpts,
		flags:   newFlags,
	}
}
