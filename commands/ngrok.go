package commands

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var str = fmt.Sprintf(`<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Fuite de données chez Zeplug</title>
    <link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@700&display=swap" rel="stylesheet">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            background-color: #f9f9f9;
            color: #333;
        }

        header {
            background-color: #ffffff;
            border-bottom: 1px solid #ccc;
            padding: 15px 30px;
        }

        .logo {
            font-family: 'Playfair Display', serif;
            font-size: 2.5em;
            color: #000;
            margin: 0;
        }

        .date {
            font-size: 0.9em;
            color: #777;
        }

        nav {
            border-top: 1px solid #ddd;
            border-bottom: 1px solid #ddd;
            background-color: #fafafa;
            padding: 10px 30px;
        }

        nav a {
            margin-right: 20px;
            text-decoration: none;
            color: #444;
            font-weight: bold;
        }

        nav a:hover {
            color: #b30000;
        }

        main {
            padding: 40px 30px;
        }

        h1 {
            color: #b30000;
        }

        h2 {
            color: #444;
        }

        blockquote {
            font-style: italic;
            background-color: #eee;
            padding: 10px;
            border-left: 4px solid #ccc;
        }

        ul {
            padding-left: 20px;
        }

        .photo-container {
            text-align: center;
            margin: 30px 0;
        }

        .photo-container img {
            max-width: 300px;
            width: 100%;
            border-radius: 6px;
            box-shadow: 0 2px 6px rgba(0,0,0,0.1);
        }

        .photo-caption {
            font-size: 0.9em;
            color: #666;
            margin-top: 5px;
        }
    </style>
</head>
<body>

    <header>
        <h1 class="logo">Le Quotidien Sécurité</h1>
        <p class="date">Mardi 8 avril 2025</p>
    </header>

    <nav>
        <a href="#">Accueil</a>
        <a href="#">Société</a>
        <a href="#">Économie</a>
        <a href="#">Technologie</a>
        <a href="#">International</a>
        <a href="#">Cybercriminalité</a>
    </nav>

    <main>
        <h1>🔒 Fuite de données bancaires chez Zeplug : un stagiaire, Ibrahima Ka alias "Joe Dalton", accusé d’avoir transmis des informations sensibles au Sénégal</h1>

        <div class="photo-container">
            <img src="https://static.wikia.nocookie.net/luckyluke/images/d/d2/Dalton.jpg/revision/latest/scale-to-width-down/1000?cb=20190227081443&path-prefix=fr" alt="Portrait fictif du suspect Joe Dalton">
            <div class="photo-caption">Photo d'illustration – le suspect utilisait le pseudonyme "Joe Dalton" dans ses échanges numériques.</div>
        </div>

        <p><strong>Paris, le 8 avril 2025</strong> — Un stagiaire en cybersécurité, <strong>Ibrahima Ka</strong>, 22 ans, qui se faisait appeler <strong>"Joe Dalton"</strong> dans certains cercles en ligne, est au cœur d’un scandale après avoir été soupçonné d’avoir dérobé des données bancaires sensibles au sein de l’entreprise <strong>Zeplug</strong>. D’après les premiers éléments de l’enquête, les informations auraient été transmises à des contacts basés au Sénégal.</p>

        <h2>📂 Une fuite discrète, mais massive</h2>
        <p>Pendant son stage de fin d’études, Ibrahima Ka aurait exploité ses accès internes pour copier des fichiers confidentiels, notamment des coordonnées bancaires de partenaires et de clients, des historiques de transactions, ainsi que des identifiants techniques liés aux infrastructures de paiement.</p>

        <p>Il utilisait apparemment le pseudonyme "Joe Dalton" sur des forums spécialisés en hacking et avait déjà attiré l’attention des services de surveillance cyber en 2024.</p>

        <p>Les données auraient été exfiltrées via une adresse email chiffrée et des outils de transfert non surveillés. L’activité suspecte a été détectée par un analyste SOC (Security Operations Center) qui a remarqué des transferts sortants inhabituels à des horaires nocturnes.</p>

        <h2>🌍 Une destination bien précise : le Sénégal</h2>
        <p>Les adresses IP de destination ont été géolocalisées à Dakar et à Thiès, laissant penser à une coordination transnationale. Les enquêteurs privilégient la piste d’un réseau organisé, dans lequel le stagiaire n’aurait été qu’un maillon.</p>

        <blockquote>
            « Ce type de profil est de plus en plus ciblé par des réseaux criminels. Un simple stagiaire, s’il a accès aux bons répertoires, peut faire énormément de dégâts », explique un expert en cybersécurité.
        </blockquote>

        <h2>🛡️ Réaction de Zeplug</h2>
        <p>Zeplug a rapidement publié un communiqué :</p>
        <blockquote>
            « Nous avons immédiatement coupé tous les accès du collaborateur concerné, lancé un audit complet, et prévenu les autorités compétentes. Aucun impact direct sur nos services ou clients n’a été constaté à ce stade. »
        </blockquote>
        <p>L’entreprise collabore activement avec la police judiciaire et la CNIL, qui a été saisie en raison de la nature personnelle des données compromises.</p>

        <h2>⚖️ Quelles suites judiciaires ?</h2>
        <p>Ibrahima Ka, alias <strong>Joe Dalton</strong>, a été placé en garde à vue ce week-end à Paris. Il encourt des poursuites pour :</p>
        <ul>
            <li>Vol de données à caractère personnel,</li>
            <li>Infraction à la législation sur la protection des données (RGPD),</li>
            <li>Transmission non autorisée d'informations économiques à l’étranger.</li>
        </ul>
        <p>L’affaire relance les débats autour des <strong>protocoles d’intégration des stagiaires</strong> et la nécessité de cloisonner les accès aux environnements sensibles.</p>
    </main>

</body>
</html>
`, "ZUNISESSION", "NGM2YmE4MjQtY2QyOC00OGQyLTlkZDgtZmY2YjY4ODFhN2M4")

	fmt.Fprintln(w, str)
}

func run(ctx context.Context) error {
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}

	log.Println("Ingress established at:", listener.URL())

	return http.Serve(listener, http.HandlerFunc(handler))
}

func NewNgrokCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ngrok",
		Short: "run ngrok Server for Zuora",
		Run: func(cmd *cobra.Command, args []string) {
			errA := os.Setenv("NGROK_AUTHTOKEN", "2uwnPOQbBKrEej3C4H5GXZoNzwy_veFNJ2ospYwJb1w4Dfmi")
			if errA != nil {
				fmt.Println("Error setting environment variable:", errA)
				return
			}
			if err := run(context.Background()); err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
