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

9. meteo hour 在浏览器中显示 done

10. vitesse加入逻辑且每个生物定死一个vitesse done
地图更新的显示请求从server端周期性主动发起变为client端周期性请求后server端响应
否则第206行的main倒计时中的updateAndSendTerrain(terr)可能会和368行的SimulateInsecte的SeDeplacer的updateAndSendTerrain(terr)发生冲突

11. dieu中修改并显示meteo done
 
12. dieu中修改并显示climat

12. 可以做一个snapshot

13. 设置每个格子密度


Objets任务

1. shit可以作为objet，作为植物的肥料

以后待办：

1. Client端转React

2. 程序什么时候结束运行（没有生物存在、到达时间上限、手动终止）

4. 参数平衡设置，现在植物长得太快了 done

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
1. 有时候昆虫停在terrain上消不掉 （被Manger的昆虫会保留在棋盘上）

2. panic报错 （并发访问和更新terrain出错) done
原因：goroutine没结束，倒计时1s结束就进入新的循环，可以使用waitgroup或者channel

您的 updateAndSendTerrain 函数中的确存在可能导致问题的并发操作。这个函数是在多个 goroutine 中被调用的，且它对全局 clients map 进行读写操作。由于 map 在 Go 中并不是并发安全的，当多个 goroutine 同时读写 map 时，就可能会导致竞态条件（race condition）和程序崩溃。
在您的代码中，clients map 被用于跟踪所有的 WebSocket 连接。您已经使用 mutex 来尝试同步对这个 map 的访问，这是正确的做法。但是，您还需要确保在调用 updateAndSendTerrain 时对这个 map 的访问是同步的。）
```

```
## blahblahblah...