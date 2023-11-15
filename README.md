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
```
1. 定义生态系统的初始状态
   首先，确定生态箱在模拟开始时的初始状态。这包括：

初始化环境：定义环境的初始参数，如温度、湿度、光照等。
创建初始生物：决定开始时生态箱中会有哪些生物（动物、植物、昆虫等）。
生物的初始属性：为每个生物设定初始属性，如年龄、能量、位置等。

2. 实现生物的基本行为
   对于生态箱中的每种生物，实现其基本行为。这可能包括：

移动：生物如何在生态箱中移动。
觅食：生物如何寻找食物。
交互：生物之间的相互作用，如捕食、竞争和繁殖。
生长和衰老：生物随时间的生长和衰老过程。

3. 模拟循环
   开发模拟循环，它将在每个时间步长（tick）中更新生态箱的状态。每个循环可能包括：

更新每个生物的状态（基于其行为和环境交互）。
更新环境状态（比如，天气变化）。
检查和处理生物的生死。

4. 数据收集和展示
   确定您想要收集哪些数据用于展示或分析。例如，您可能想要跟踪特定生物的数量、环境参数的变化等。实现将这些数据从模拟后端发送到前端的逻辑。

5. 前端集成
   在前端（网页客户端），根据从服务器接收的数据更新展示。这可能包括：

显示生态箱的当前状态。
提供用户与模拟交互的界面（如果适用）。
可视化数据（如图表或地图）。

6. 测试和调试
   开始简单，逐渐增加复杂性。定期测试模拟以确保它按预期工作，并调整参数和行为以获得更逼真的结果。

7. 性能考虑
   根据需要优化性能。复杂的模拟可能会很快消耗资源，特别是当模拟大量生物或复杂交互时。

8. 用户交互（可选）
   如果您想让用户能够影响模拟，考虑添加这些功能，例如添加或移除生物、改变环境条件等。

9. 文档和用户指南
   随着项目的发展，保持更新文档和用户指南，尤其是如果项目是为了教育目的或将被其他开发者使用。
```

## Important提示 
```
1. 为什么使用WebSocket而不是Http通信模式
- 可以让Server自发发消息给Client，用于定期检查
- WebSocket: 提供了全双工通信通道，允许数据在客户端和服务器之间双向流动。一旦建立了WebSocket连接，客户端和服务器就可以随时互相发送数据。
- 适用于需要实时通信的应用，如在线聊天、实时游戏、实时数据更新等。建立连接后，保持连接开放，减少了因多次握手导致的延迟和开销。
- 一旦建立，连接就会保持开放，直到客户端或服务器明确关闭为止。
- 初始握手包含头部信息，但之后的数据传输不再需要重复发送头部信息。

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

```
## blahblahblah...