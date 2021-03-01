// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package template_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gardener/component-cli/pkg/template"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Template Test Suite")
}

var _ = Describe("Template", func() {

	Context("Parse Arguments", func() {

		It("should parse one argument after a '--'", func() {
			opts := template.Options{}
			Expect(opts.Complete([]string{"--", "MY_VAR=test"})).To(Succeed())
			Expect(opts.Vars).To(HaveKeyWithValue("MY_VAR", "test"))
		})

		It("should parse no argument if no '--' separator is provided", func() {
			opts := template.Options{}
			Expect(opts.Complete([]string{"MY_VAR=test"})).To(Succeed())
			Expect(opts.Vars).To(HaveLen(0))
		})

		It("should parse multiple values after a '--'", func() {
			opts := template.Options{}
			Expect(opts.Complete([]string{"--", "MY_VAR=test", "myOtherVar=true"})).To(Succeed())
			Expect(opts.Vars).To(HaveKeyWithValue("MY_VAR", "test"))
			Expect(opts.Vars).To(HaveKeyWithValue("myOtherVar", "true"))
		})

	})

	Context("Template", func() {
		It("should template with a single value", func() {
			s := "my ${MY_VAR}"
			opts := template.Options{}
			opts.Vars = map[string]string{
				"MY_VAR": "test",
			}
			res, err := opts.Template(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my test"))
		})

		It("should template multiple value", func() {
			s := "my ${MY_VAR} ${my_second_var}"
			opts := template.Options{}
			opts.Vars = map[string]string{
				"MY_VAR":        "test",
				"my_second_var": "testvalue",
			}
			res, err := opts.Template(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my test testvalue"))
		})

		It("should use an empty string if no value is provided", func() {
			s := "my ${MY_VAR}"
			opts := template.Options{}
			opts.Vars = map[string]string{}
			res, err := opts.Template(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my "))
		})

	})

})
