// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by "mdtogo"; DO NOT EDIT.
package producer

var CompositepkgsGuide = `
{{% pageinfo color="warning" %}}

#### Notice: Composite packages feature support is in alpha phase

{{% /pageinfo %}}

This guide walks you through an example to create, parameterize and publish a kpt
package which has nested subpackages in it. A kpt package is a directory of resource
configs with valid ` + "`" + `Kptfile` + "`" + ` in it. A composite package is a ` + "`" + `kpt` + "`" + ` package with 1
or more subpackages in it.

Principles:

1. Each kpt package is an independent building block and
   should contain resources(ex: setter definitions) of its own.
2. If a package is present in the directory tree of parent package,
   the configs of that package are out of scope for the actions performed
   on the parent package.
3. ` + "`" + `--recurse-subpackages(-R)` + "`" + ` flag can be leveraged to run a command
   recursively on all the nested subpackages of a root package.

## Steps

1. [Create a composite package](#create-a-composite-package)
2. [Add setters and substitutions](#add-setters-and-substitutions)
3. [Publish the package](#publish-the-package)

## Create a composite package

Create a composite package directory structure and initialize ` + "`" + `kpt` + "`" + ` packages.

  mkdir hello-composite-pkg
  kpt pkg init hello-composite-pkg
  
  mkdir hello-composite-pkg/hello-subpkg
  kpt pkg init hello-composite-pkg/hello-subpkg
  
  # this is a subdir and not a package
  mkdir hello-composite-pkg/hello-subpkg/hello-dir
  
  mkdir hello-composite-pkg/hello-subpkg/hello-nestedpkg
  kpt pkg init hello-composite-pkg/hello-subpkg/hello-nestedpkg

Add resource files(.yaml) to the directories. You may use(copy/paste) the resource files from
[hello-composite-pkg] to respective directories on local and delete the ` + "`" + `$kpt-set` + "`" + ` comments.

## Add setters and substitutions

### Create setters

[Setters] provide a solution for template-free setting of field values through
package metadata (OpenAPI). Setters will be invoked by package consumers to
programmatically modify the configuration using ` + "`" + `kpt cfg set` + "`" + ` to [set] values.
Create ` + "`" + `namespace` + "`" + ` setter in all the packages.

  kpt cfg create-setter hello-composite-pkg/ namespace YOURSPACE -R --required

Output:

  hello-composite-pkg/
  created setter "namespace"
  
  hello-composite-pkg/hello-subpkg/
  created setter "namespace"
  
  hello-composite-pkg/hello-subpkg/hello-nestedpkg/
  created setter "namespace"

Similarly create a setter with name ` + "`" + `gcloud.core.project` + "`" + `. If the package consumer
has ` + "`" + `gcloud` + "`" + ` set up on local, they can observe that the value of the setter
` + "`" + `gcloud.core.project` + "`" + ` will be set automatically when the package is fetched.
` + "`" + `gcloud` + "`" + ` config setters are automatically set deriving the values from the
output of ` + "`" + `gcloud config list` + "`" + ` command, when the package is fetched using [kpt pkg get].

  kpt cfg create-setter hello-composite-pkg/ gcloud.core.project PROJECT_ID -R

### Create substitutions

[Substitutions] provide a solution for template-free substitution of field values
built on top of setters. They enable substituting values into part of a field,
including combining multiple setters into a single value.

Substitutions may be invoked to programmatically modify the configuration using
` + "`" + `kpt cfg set` + "`" + ` to substitute values which are derived from the setter.

  kpt cfg create-subst hello-composite-pkg/ image-tag \
  --field-value gcr.io/kpt-dev/helloworld-gke:0.1.0 \
  --pattern gcr.io/kpt-dev/\${image}:\${tag} -R

Output:

  hello-composite-pkg/
  unable to find setter with name image, creating new setter with value helloworld-gke
  unable to find setter with name tag, creating new setter with value 0.1.0
  created substitution "image-tag"
  
  hello-composite-pkg/hello-subpkg/
  unable to find setter with name image, creating new setter with value helloworld-gke
  unable to find setter with name tag, creating new setter with value 0.1.0
  created substitution "image-tag"
  
  hello-composite-pkg/hello-subpkg/hello-nestedpkg/
  unable to find setter with name image, creating new setter with value helloworld-gke
  unable to find setter with name tag, creating new setter with value 0.1.0
  created substitution "image-tag"

### List and verify setters/substitutions

Use list-setters command to verify that the setters and substitutions are created as expected

  kpt cfg list-setters hello-composite-pkg/ --include-subst

Output:

  hello-composite-pkg/
           NAME                VALUE        SET BY   DESCRIPTION   COUNT   REQUIRED
    gcloud.core.project   YOUR_PROJECT_ID                          1       No
    image                 helloworld-gke                           1       No
    namespace             YOURSPACE                                1       Yes
    tag                   0.1.0                                    1       No
    SUBSTITUTION              PATTERN               REFERENCES
    image-tag      gcr.io/kpt-dev/${image}:${tag}   [image,tag]
  
  hello-composite-pkg/hello-subpkg/hello-nestedpkg/
           NAME                VALUE        SET BY   DESCRIPTION   COUNT   REQUIRED
    gcloud.core.project   YOUR_PROJECT_ID                          1       No
    image                 helloworld-gke                           1       No
    namespace             YOURSPACE                                1       Yes
    tag                   0.1.0                                    1       No`