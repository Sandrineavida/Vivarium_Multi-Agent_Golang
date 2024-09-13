# Projet Vivarium

## Auteurs

- Zhenyang Xu - zhenyang.xu@etu.utc.fr
- Hudie Sun - hudie.sun@etu.utc.fr
- Jinxing Lai - jinxing.lai@etu.utc.fr


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
├── static
│   └── index.html                     # Affichage dans navigateur
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

### Cloner le Projet

Pour cloner ce projet, veuillez utiliser la branche `LJX` en exécutant la commande suivante :

```bash
git clone -b LJX https://gitlab.utc.fr/xuzhenya/projet-ia04-vivarium
```

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

## Captures d'écran du Projet

Le projet Vivarium est visualisé à travers plusieurs GIFs illustrant les divers aspects de notre simulation. Voici un aperçu des fonctionnalités clés et des interactions au sein de l'écosystème virtuel.

### Interface Utilisateur et Interaction HTML

**Démonstration du HTML et du WebSocket**  
*Description*: Affichage HTML simplifié montrant l'ensemble de l'écosystème avec des options pour ajouter des insectes, des plantes et modifier le climat.  

![Démonstration](/Demonstration/vvtm-html.gif)

### Interactions et Comportements

**Démonstration d’Interaction**  
*Description*: Interaction entre les différents organismes au sein de l'écosystème.  

<p align="center">
  <img src="/Demonstration/decorations_sur_la_tete.jpg" alt="Décoration sur la tête" width="300"/>&nbsp;&nbsp;&nbsp;&nbsp;
  <img src="/Demonstration/demonstration_d'interaction1.gif" alt="Démonstration1" width="175"/>&nbsp;&nbsp;&nbsp;&nbsp;
  <img src="/Demonstration/demonstration_d'interaction2.gif" alt="Démonstration2" width="230"/>
</p>


### Contrôle de Simulation

**Démonstration du PauseSignal**  
*Description*: Fonctionnalité permettant de mettre en pause et de reprendre la simulation.  

<div align="center">
  <img src="/Demonstration/demonstration_du_PauseSignal.gif" alt="Démonstration" />
</div>

### Système Climatique

**Démonstration du Système Climatique**  

*Brouillard*:

<div align="center">
  <img src="/Demonstration/meteo_brouillard.gif" alt="Démonstration Brouillard" />
</div>

*Pluie*: 

<div align="center">
  <img src="/Demonstration/meteo_pluie.gif" alt="Démonstration Pluie" />
</div>

*Incendie*: 

<div align="center">
  <img src="/Demonstration/meteo_incendie.gif" alt="Démonstration Incendie" />
</div>

*Tonnerre*: 

<div align="center">
  <img src="/Demonstration/meteo_tonerre.gif" alt="Démonstration Tonnerre" />
</div>

### Écosystème Vivarium

**Démonstration du Vivarium**  
*Description*: Comportements et reproduction des espèces dans l'écosystème.  

Dans ce GIF de démonstration : Comportements et Reproduction des Espèces
1. **Araignée Sauteuse recherche un partenaire de reproduction et se bat avec un individu du même sexe.**
2. **Lombric en processus de reproduction.**
3. **Reproduction végétale.**
4. **Petit Serpent se nourrit par faim.**
5. **Déplacement d'insectes.**
6. **Les insectes meurent de vieillesse ou d'être mangés, et les plantes meurent d'être mangées...**

<div align="center">
  <img src="/Demonstration/demonstration_du_Vivaraium.gif" alt="Démonstration du Vivarium" />
</div>

### Intervention Divine

**Démonstration de Comportement de Dieu**  
*Description*: Interaction et impact de l'intervention 'divine' sur l'écosystème.  

<div align="center">
  <img src="/Demonstration/demonstation_de_comportement_dieu.gif" alt="Démonstration de Comportement de Dieu" />
</div>

## Remerciements

Un grand merci à tous ceux qui ont contribué au projet, en particulier au professeur Sylvain Lagrue pour son soutien technique.
