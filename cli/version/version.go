/*
	Copyright 2023 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package version

import (
	"github.com/loopholelabs/cmdutils/pkg/version"
	"github.com/loopholelabs/scale/cli/internal/config"
)

var (
	// GitCommit is filled in at build time and contains the last git commit hash when this application was built
	GitCommit = ""

	// GoVersion is filled in at build time and contains the golang version upon which this application was built
	GoVersion = ""

	// Platform is filled in at build time and contains the platform upon which this application was built
	Platform = ""

	// Version is filled in at build time and contains the official release version of this application
	Version = ""

	// BuildDate is filled in at build time and contains the date when this application was build
	BuildDate = ""
)

var V = version.New[*config.Config](GitCommit, GoVersion, Platform, Version, BuildDate)
