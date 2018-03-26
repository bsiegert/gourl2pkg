/*-
* Copyright (c) 2018
*  Benny Siegert <bsiegert@gmail.com>
*
* Provided that these terms and disclaimer and all copyright notices
* are retained or reproduced in an accompanying document, permission
* is granted to deal in this work without restriction, including un-
* limited rights to use, publicly perform, distribute, sell, modify,
* merge, give away, or sublicence.
*
* This work is provided "AS IS" and WITHOUT WARRANTY of any kind, to
* the utmost extent permitted by applicable law, neither express nor
* implied; without malicious intent or gross negligence. In no event
* may a licensor, author or contributor be held liable for indirect,
* direct, other damage, loss, or other issues arising in any way out
* of dealing in the work, even if advised of the possibility of such
* damage or existence of a defect, except proven that it results out
* of said person's immediate fault when using the work as intended.
 */

package pkgsrc

import "text/template"

const makefileTemplate = `# $NetBSD$

DISTNAME=               {{.Distname}}
CATEGORIES=             {{for range .Categories}}{{.}} {{end}}
MASTER_SITES=           {{.MasterSites}}
#GITHUB_PROJECT=
#GITHUB_TAG=


MAINTAINER=             pkgsrc-users@NetBSD.org
HOMEPAGE=               TODO: add homepage
COMMENT=                TODO: add comment
#LICENSE=                # TODO

GO_SRCPATH=             {{.GoSrcpath}}
#GO_DIST_BASE=           ${GITHUB_PROJECT}-${GITHUB_TAG}*

CHECK_RELRO_SKIP+=      # TODO: add any binaries here

{{.AllDependencies}}
.include "../../lang/go/go-package.mk"
.include "../../mk/bsd.pkg.mk"
`

var makefileTmpl = template.Must(template.New("Makefile").Parse(makefileTemplate))
