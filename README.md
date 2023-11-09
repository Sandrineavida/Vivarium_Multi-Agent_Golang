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
│   ├── insectes.go                    # Définit l'énumération des types d'insectes
│   └── plantes.go                     # Définit l'énumération des types de plantes
├── climat                             # Code relatif au "Climat"
│   └── climat.go                      # Définit la structure Climat et ses méthodes
├── dieu                               # Code relatif au "Dieu"
│   └── dieu.go                        # Définit la structure Dieu et ses méthodes
├── environnement                      # Code pour "Environment" et "Terrain"
│   ├── environnement.go               # Définit la structure Environment et ses méthodes
│   └── terrain.go                     # Définit la structure Terrain et ses méthodes
└── utils                              # Répertoire pour les fonctions utilitaires et le code commun
    └── utils.go                       # Fonctions utilitaires et code commun
```

## blahblahblah...