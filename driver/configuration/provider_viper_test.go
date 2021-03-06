package configuration

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/rs/cors"

	"github.com/ory/hydra/x"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func setEnv(key, value string) func(t *testing.T) {
	return func(t *testing.T) {
		require.NoError(t, os.Setenv(key, value))
	}
}

func TestSubjectTypesSupported(t *testing.T) {
	p := NewViperProvider(logrus.New(), false, nil)
	viper.Set(ViperKeySubjectIdentifierAlgorithmSalt, "00000000")
	for k, tc := range []struct {
		d string
		p func(t *testing.T)
		e []string
		c func(t *testing.T)
	}{
		{
			d: "Load legacy environment variable in legacy format",
			p: setEnv(strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)), "public,pairwise,foobar"),
			c: setEnv(strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)), ""),
			e: []string{"public", "pairwise"},
		},
		{
			d: "Load legacy environment variable in legacy format",
			p: func(t *testing.T) {
				setEnv(strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)), "public,pairwise,foobar")(t)
				setEnv(strings.ToUpper(strings.Replace(ViperKeyAccessTokenStrategy, ".", "_", -1)), "jwt")(t)
			},
			c: setEnv(strings.ToUpper(strings.Replace(ViperKeySubjectTypesSupported, ".", "_", -1)), ""),
			e: []string{"public"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			tc.p(t)
			assert.EqualValues(t, tc.e, p.SubjectTypesSupported())
			tc.c(t)
		})
	}
}

func TestWellKnownKeysUnique(t *testing.T) {
	p := NewViperProvider(logrus.New(), false, nil)
	assert.EqualValues(t, []string{x.OAuth2JWTKeyName, x.OpenIDConnectKeyName}, p.WellKnownKeys(x.OAuth2JWTKeyName, x.OpenIDConnectKeyName, x.OpenIDConnectKeyName))
}

func TestCORSOptions(t *testing.T) {
	p := NewViperProvider(logrus.New(), false, nil)
	viper.Set("serve.public.cors.enabled", true)

	assert.EqualValues(t, cors.Options{
		AllowedOrigins:     []string{},
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:     []string{"Authorization"},
		ExposedHeaders:     []string{"Content-Type"},
		AllowCredentials:   true,
		OptionsPassthrough: false,
		MaxAge:             0,
		Debug:              false,
	}, p.CORSOptions("public"))
}
