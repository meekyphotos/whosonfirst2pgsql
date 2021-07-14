package entrypoints

import (
	"github.com/meekyphotos/whosonfirst2pgsql/pkg/commands"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockRunner struct{}

func (pr MockRunner) Run(config *commands.Config) error {
	return nil
}

func TestArgumentsAreParsedCorrectly(t *testing.T) {
	underTest := commands.Loader{
		Runner: MockRunner{},
	}

	app := InitApp(&underTest)
	err := app.Run([]string{"app", "-c", "-d", "db_test", "hello"})

	assert.Nil(t, err)

	assert.NotEmpty(t, underTest.Config.File)
	assert.Equal(t, true, underTest.Config.Create)
}

func TestAppendModeSetCreateToFalse(t *testing.T) {
	underTest := commands.Loader{
		Runner: MockRunner{},
	}

	app := InitApp(&underTest)
	err := app.Run([]string{"app", "-a", "-d", "db_test", "hello"})

	assert.Nil(t, err)

	assert.NotEmpty(t, underTest.Config.File)
	assert.Equal(t, false, underTest.Config.Create)
}

type MockPasswordReader struct{}

func (pr MockPasswordReader) ReadPassword() (string, error) {
	return "promptedPassword", nil
}

func TestDatabaseOptions(t *testing.T) {
	underTest := commands.Loader{
		Runner:           MockRunner{},
		PasswordProvider: MockPasswordReader{},
	}

	app := InitApp(&underTest)
	err := app.Run([]string{"app", "-a", "-d", "db_test", "-U", "meeky", "-W", "-P", "13245", "hello"})
	assert.Nil(t, err)

	assert.Equal(t, "db_test", underTest.Config.DbName)
	assert.Equal(t, "meeky", underTest.Config.UserName)
	assert.Equal(t, "promptedPassword", underTest.Config.Password)
	assert.Equal(t, "localhost", underTest.Config.Host)
	assert.Equal(t, 13245, underTest.Config.Port)

}

func TestOutputFormatLatLngDisablesGeom(t *testing.T) {
	underTest := commands.Loader{
		Runner:           MockRunner{},
		PasswordProvider: MockPasswordReader{},
	}

	app := InitApp(&underTest)
	err := app.Run([]string{"app", "-d", "db_test", "--latlong", "hello"})
	assert.Nil(t, err)

	assert.Equal(t, true, underTest.Config.UseLatLng)
	assert.Equal(t, false, underTest.Config.UseGeom)

}

func TestOutputFormat(t *testing.T) {
	underTest := commands.Loader{
		Runner:           MockRunner{},
		PasswordProvider: MockPasswordReader{},
	}

	app := InitApp(&underTest)
	err := app.Run([]string{"app", "-d", "db_test", "-p", "tbl", "-k", "hello"})
	assert.Nil(t, err)

	assert.Equal(t, true, underTest.Config.UseGeom)
	assert.Equal(t, "tbl", underTest.Config.TableName)
	assert.Equal(t, true, underTest.Config.InclKeyValues)
	assert.Equal(t, true, underTest.Config.ExcludeColumnFromKeyValues)
	assert.Equal(t, false, underTest.Config.UseJson)
}

func TestOutputFormatHStoreAll(t *testing.T) {
	underTest := commands.Loader{
		Runner:           MockRunner{},
		PasswordProvider: MockPasswordReader{},
	}

	app := InitApp(&underTest)
	err := app.Run([]string{"app", "-d", "db_test", "-j", "hello"})
	assert.Nil(t, err)

	assert.Equal(t, true, underTest.Config.UseGeom)
	assert.Equal(t, true, underTest.Config.InclKeyValues)
	assert.Equal(t, false, underTest.Config.ExcludeColumnFromKeyValues)
	assert.Equal(t, false, underTest.Config.UseJson)
	assert.Equal(t, false, underTest.Config.MatchOnly)
}

func TestOutputFormatJsonbAll(t *testing.T) {
	underTest := commands.Loader{
		Runner:           MockRunner{},
		PasswordProvider: MockPasswordReader{},
	}

	app := InitApp(&underTest)
	err := app.Run([]string{"app", "-d", "db_test", "--json-all", "--match-only", "--schema", "custom", "hello"})
	assert.Nil(t, err)

	assert.Equal(t, true, underTest.Config.UseGeom)
	assert.Equal(t, true, underTest.Config.InclKeyValues)
	assert.Equal(t, false, underTest.Config.ExcludeColumnFromKeyValues)
	assert.Equal(t, true, underTest.Config.UseJson)
	assert.Equal(t, true, underTest.Config.MatchOnly)
	assert.Equal(t, "custom", underTest.Config.Schema)
}
