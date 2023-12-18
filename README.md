# Projet Vivarium

## Auteurs

- Zhenyang Xu - zhenyang.xu@etu.utc.fr
- Hudie Sun - hudie.sun@etu.utc.fr
- Jinxing Lai - jinxing.lai@etu.utc.fr
- Noe Redouin - noe.redouin@etu.utc.fr

## À propos du projet

Le projet Vivarium est une simulation d'écosystème dynamique où divers organismes, 
tels que des insectes et des plantes, interagissent dans un environnement virtuel. 
Ce projet utilise le langage Go pour la logique du serveur et Ebiten pour le rendu graphique, 
offrant une visualisation en temps réel de l'écologie d'un territoire virtuel.

## Structure du Projet

```plaintext
/vivarium                              # Répertoire racine du projet
├── go.mod                             # Fichier du module Go
├── server.go                          # Serveur pour simuler et pour gérer la communication avec ebiten et hmtl
├── index.html                         # Affichage dans navigateur
├── ebiten                             # Affichage par Ebiten
│   ├── assets                         
│   │   ├── images                     # Stockage des ressoureces graphiques
│   │   │    ├── somepic.png
│   │   │    ├── somepic.go
│   │   │    ...
│   │   └── file2byteslice.go          # Conversion du format PNG en BYTESLICE
│   ├── sprites.go                     # Définit de la logique de rendu pour chaque espèce
│   └── ebiton.go                      # Point d'entrée du programme et gérer le moteur de rendu Ebiten
├── organisme                          # Code pour "Organisme" et ses sous-classes
│   ├── organisme.go                   # Définit l'interface Organisme et la structure BaseOrganisme
│   ├── insectes.go                    # Définit la structure Insecte et ses méthodes
│   └── plantes.go                     # Définit la structure Plante et ses méthodes
├── enums                              # Définitions des types énumérés
│   ├── sexe.go                        # Définit l'énumération Sexe
│   ├── meteo.go                       # Définit l'énumération Meteo
│   └── espece.go                      # Définit l'énumération des types d'insectes et de plantes
├── climat                             # Code relatif au "Climat"
│   └── climat.go                      # Définit la structure Climat et ses méthodes
├── environnement                      # Code pour "Environnement"
│   └── environnement.go               # Définit la structure Environment et ses méthodes
├── terrain                            # Code pour "Terrain"
│   └── terrain.go                     # Définit la structure Terrain et ses méthodes
└── utils                              # Répertoire pour les fonctions utilitaires et le code commun
    └── utils.go                       # Fonctions utilitaires et code commun
```

## Installation

### Prérequis

- Go version 1.x
- Ebiten

### Installation des dépendances

```shell
go get -u github.com/gorilla/websocket
go get github.com/hajimehoshi/ebiten/v2
```

## Lancement

Pour démarrer le serveur et lancer l'affichage Ebiten:

```shell
go run ebiten/ebiten.go
```

Pour afficher de simulation HTML ouvert et formulaire de fonction Dieu
```
http://localhost:8000/
```

## Remerciements

Un grand merci à tous ceux qui ont contribué au projet, en particulier au professeur Sylvain Lagrue pour son soutien technique.