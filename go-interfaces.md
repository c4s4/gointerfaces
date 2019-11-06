---
title:      Les interfaces du Go
author:     Michel Casabianca
date:       2014-10-28
updated:    UPDATE
categories: [articles]
tags:       [golang]
id:         go-interfaces
email:      casa@sweetohm.net
lang:       fr
toc:        no
---

En assistant à la dotGo, où le buzzword était clairement *l'interface*, je me suis demandé où l'on pouvait se procurer la liste de toutes les interfaces définies dans le langage. J'ai cherché et n'ai trouvé cette information nulle part.

Je me suis donc décidé à écrire un petit programme qui :

- Télécharge le tarball des sources d'une version donnée.
- Parse les fichiers sources pour en extraire les interfaces ainsi que le numéro de ligne où elles sont définies.
- Affiche sur la console la liste de ces interfaces sous la forme d'un tableau markdown.

<!--more-->

Le projet se trouve sur Github : <https://github.com/c4s4/gointerfaces>.

Voici le résultat :

INTERFACES

On pourra trouver une discussion sur ces interfaces dans l'article suivant (en anglais) : <http://mwholt.blogspot.fr/2014/08/maximizing-use-of-interfaces-in-go.html>.

*Enjoy!*
