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

1. 添加动植物也要加锁 done

2. 植物逻辑加入server done
原因1：checkEtat逻辑写反了
原因2：新生植物坐标逻辑修改，根据半径随机生成新植物，但如果植物在边界，新生植物坐标得限定范围

4. 动物打架 done

5. 性冷淡问题解决 done 
原因：生孩子数量和养胃年纪构造函数中填反了

6. 交配函数报错 done 
原因：找异性或同性target交配的结果可能为nil，需要if判断

7. 清理一下print，print断点测试太多了

8. 添加昆虫的时候新增一个限制，蜗牛和蟋蟀性别定死只能单性别，其他昆虫不能选择单性别，只能选择男女 done

我做不动，法国哥做一下：

1. Client端转React

2. 程序什么时候结束运行（没有生物存在、到达时间上限、手动终止）

3. climat的所有逻辑

4. 参数平衡设置，现在植物长得太快了

5. 移植2d游戏引擎
```

## 结构：
```
对于每一个goroutine对应的agent，三个动作
Percetpt
Deliberate
Act
```

## 霸哥：
```
1. 有时候昆虫停在terrain上消不掉

2. panic报错
```

```
## blahblahblah...