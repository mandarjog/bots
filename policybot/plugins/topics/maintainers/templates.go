// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package maintainers

// TODO: in the JavaScript, remove org=istio and assume the same org the page was initially rendered with
var maintainerTemplate = `
{{ define "content" }}

<p>
These kind folks are responsible for specific areas of the Istio product, guiding its development
and maintaining its code base.
</p>

<table>
    <thead>
    <tr>
        <th>Avatar</th>
        <th>Login</th>
        <th>Name</th>
        <th>Company</th>
        <th>Emeritus</th>
    </tr>
    </thead>
    <tbody>
        {{ range . }}
            <tr>
                <td><img style='width: 30px' src='{{ .AvatarURL }}'/></td>
                <td>{{ .Login }}</td>
                <td>{{ .Name }}</td>
                <td>{{ .Company }}</td>
                <td>{{ .Emeritus }}</td>
            </tr>
        {{ end }}
    </tbody>
</table>
{{ end }}
`
