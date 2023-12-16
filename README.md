## Structure du Projet

```plaintext
/vivarium                              # Répertoire racine du projet
├── go.mod                             # Fichier du module Go
├── server.go                          # Serveur pour simuler et pour gérer la communication avec ebiten et hmtl
├── index.html                         # Affichage dans navigateur
├── ebiten                             # Affichage par Ebiten
│   ├── assets                         
|   |   ├── images                     # Stockage des ressoureces graphiques
|   |   |    ├── somepic.png
|   |   |    ├── somepic.go
|   |   |    ...
|   |   └── file2byteslice.go          # Conversion du format PNG en BYTESLICE
│   ├── sprites.go                     # Définit de la logique de rendu pour chaque espèce
|   └── ebiton.go                      # Point d'entrée du programme et gérer le moteur de rendu Ebiten
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

## Server et Client
Download WebSocket:
```
go get -u github.com/gorilla/websocket
go get github.com/hajimehoshi/ebiten/v2
```

Launch Server:
```
go run server.go
go run animation.go
```

## Étapes d'intégration


## Important提示 
```
Pourquoi utiliser WebSocket au lieu du mode de communication HTTP
- Le serveur peut être autorisé à envoyer spontanément des messages au client pour une vérification régulière
- WebSocket : Fournit un canal de communication full-duplex, permettant aux données de circuler dans les deux sens entre le client et le serveur. Une fois la connexion WebSocket établie, le client et le serveur peuvent s'envoyer des données à tout moment.
- Convient aux applications nécessitant une communication en temps réel, telles que le chat en ligne, les jeux en temps réel, les mises à jour de données en temps réel, etc. Une fois la connexion établie, le maintien de la connexion ouverte réduit le délai et la surcharge provoqués par plusieurs poignées de main.
- Une fois établie, la connexion reste ouverte jusqu'à ce qu'elle soit explicitement fermée par le client ou le serveur.
- La prise de contact initiale contient des informations d'en-tête, mais les transmissions de données ultérieures n'ont plus besoin d'envoyer des informations d'en-tête à plusieurs reprises.

```

## blahblahblah...