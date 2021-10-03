package ssh

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/ahmadwaleed/chore/pkg/config"
	"github.com/ahmadwaleed/chore/pkg/executer"
	"github.com/ahmadwaleed/chore/pkg/ssh"
	"github.com/ahmadwaleed/choreui/app/database/model"
	"golang.org/x/sync/errgroup"
)

type Runner struct {
	Parallel bool
	privKeys []*os.File
}

func (r *Runner) RunTask(task model.Task, callback executer.OnRunCallback) error {
	var hosts []ssh.Config
	for _, s := range task.Servers {
		f, err := ioutil.TempFile("", "id_rda_")
		if err != nil {
			return err
		}
		f.WriteString(s.SSHPrivateKey)

		r.privKeys = append(r.privKeys, f)
		hosts = append(hosts, ssh.Config{
			User:   s.User,
			Host:   s.IP,
			Port:   strconv.Itoa(s.Port),
			RSA_ID: f.Name(),
		})
	}

	sshTask := config.Task{
		Name:     task.Name,
		Env:      config.EnvVar(task.EnvVar()),
		Commands: strings.Split(task.Script, "\n"),
		Hosts:    hosts,
	}

	exec := executer.New("ssh")
	if r.Parallel {
		exec = executer.New("parallel")
	}

	return exec.Run(sshTask, callback)
}

func (r *Runner) Close() error {
	g := new(errgroup.Group)
	for i := range r.privKeys {
		f := r.privKeys[i]
		g.Go(func() error {
			return os.Remove(f.Name())
		})
	}

	return g.Wait()
}
