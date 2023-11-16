## Structure du Projet

```plaintext
/vivarium                              # Répertoire racine du projet
├── go.mod                             # Fichier du module Go
├── main.go                            # Point d'entrée du programme
├── organisme                          # Code pour "Organisme" et ses sous-classes
│   ├── organisme.go                   # Définit l'interface Organisme et la structure BaseOrganisme
│   ├── insectes.go                    # Définit la structure Insecte et ses méthodes
│   └── plantes.go                     # Définit la structure Plante et ses méthodes
├── enums                              # Définitions des types énumérés
│   ├── sexe.go                        # Définit l'énumération Sexe
│   ├── meteo.go                       # Définit l'énumération Meteo
│   ├── mode_reproduction.go           # Définit l'énumération ModeReproduction
│   ├── espece.go                    # Définit l'énumération des types d'insectes et de plantes
├── climat                             # Code relatif au "Climat"
│   └── climat.go                      # Définit la structure Climat et ses méthodes
├── dieu                               # Code relatif au "Dieu"
│   └── dieu.go                        # Définit la structure Dieu et ses méthodes
├── environnement                      # Code pour "Environnement"
│   ├── environnement.go               # Définit la structure Environment et ses méthodes
├── terrain                            # Code pour "Terrain"
│   └── terrain.go                     # Définit la structure Terrain et ses méthodes
└── utils                              # Répertoire pour les fonctions utilitaires et le code commun
    └── utils.go                       # Fonctions utilitaires et code commun
```

## Server et Client
Download WebSocket:
```
go get -u github.com/gorilla/websocket
```

Launch Server:
```
go run server.go
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

## 目前代办：
```
接下来完成这些逻辑：

逻辑1：Done
NiveauFaim上限应为10 下限应为0
Energy上限应为10 下限应为0

NiveauFaim=10或Energy=0给我去死

逻辑2：Done
每次移动，能量-1，饥饿等级+1

逻辑3：Done
植物，昆虫生长，衰老后给我去死
为每个物种新增MaxAge和AgeRate

逻辑4：
新增上帝按钮，在选定坐标添加生物和他的参数

逻辑5：
每种植物的rayon, vitesseDeCroissance，modeReproduction需要根据Espece自动添加
每种昆虫的rayon，periodReproduire需要根据Espece自动添加

html中espece不应该让dieu填写，而是应该搞成选项框

```
## blahblahblah...