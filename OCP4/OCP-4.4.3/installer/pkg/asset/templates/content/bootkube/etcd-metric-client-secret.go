package bootkube

import (
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/templates/content"
)

const (
	etcdMetricClientSecretFileName = "etcd-metric-client-secret.yaml.template"
)

var _ asset.WritableAsset = (*EtcdMetricClientSecret)(nil)

// EtcdMetricClientSecret is an asset for the etcd client cert for metrics
type EtcdMetricClientSecret struct {
	FileList []*asset.File
}

// Dependencies returns all of the dependencies directly needed by the asset
func (t *EtcdMetricClientSecret) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Name returns the human-friendly name of the asset.
func (t *EtcdMetricClientSecret) Name() string {
	return "EtcdMetricClientSecret"
}

// Generate generates the actual files by this asset
func (t *EtcdMetricClientSecret) Generate(parents asset.Parents) error {
	fileName := etcdMetricClientSecretFileName
	data, err := content.GetBootkubeTemplate(fileName)
	if err != nil {
		return err
	}
	t.FileList = []*asset.File{
		{
			Filename: filepath.Join(content.TemplateDir, fileName),
			Data:     []byte(data),
		},
	}
	return nil
}

// Files returns the files generated by the asset.
func (t *EtcdMetricClientSecret) Files() []*asset.File {
	return t.FileList
}

// Load returns the asset from disk.
func (t *EtcdMetricClientSecret) Load(f asset.FileFetcher) (bool, error) {
	file, err := f.FetchByName(filepath.Join(content.TemplateDir, etcdMetricClientSecretFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	t.FileList = []*asset.File{file}
	return true, nil
}
