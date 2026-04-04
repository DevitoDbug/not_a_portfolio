package api

import (
	"net/http"

	"github.com/DevitoDbug/portfolio/internals/web/pages"
)

func (a *Api) ResumeHandler(w http.ResponseWriter, r *http.Request) {
	experiences := []pages.Experience{
		{
			JobName:    "Shupav Group Ltd",
			TimePeriod: "from Jan 2025",
			Links: []pages.Link{
				"https://app.firmstamp.com/",
				"https://ongeri.shupavgroup.com/",
			},
			Description: []pages.Paragraph{
				"Designed and deployed production-grade Go REST APIs using clean architecture, with comprehensive unit and integration tests (Testify) across a  multi-tenant system, alongside building Next.js and Laravel applications for other products within the organization.",
				"Owned the full infrastructure setup — CI/CD pipelines (GitHub Actions), Docker containerization, CapRover deployment, and AWS cloud infrastructure for the organization.",
			},
		},
		{
			JobName:    "Sphere Labs | Remote — Hong Kong",
			TimePeriod: "Sep 2025 to Jan 2026",
			Links: []pages.Link{
				"https://www.simbank.com/",
				"https://app.sphere-id.com/",
				"https://www.liquidroyalty.com/",
			},
			Description: []pages.Paragraph{
				"Built and maintained React applications for an international team, focusing on rewriting poor UX flows and delivering new platforms with clean, responsive interfaces.",
				"Integrated DAP and REST APIs, handling client-server data flow and collaborating remotely across flexible sprints.",
			},
		},
		{
			JobName:    "Mobipine",
			TimePeriod: "June 2024 – May 2025",
			Links:      []pages.Link{},
			Description: []pages.Paragraph{
				"Built custom Odoo modules in Python across internal business systems, including integrating ZKTeco biometric hardware directly with the Odoo HR module for automated attendance tracking, alongside delivering e-commerce platforms and custom applications in Laravel.",
				"Designed a tax automation platform processing 1,000+ daily transactions with QuickBooks and Zoho Books integrations, built API endpoints handling 10K+ records monthly, and managed CI/CD pipelines via GitHub Actions.",
			},
		},
		{
			JobName:    "Sakah",
			TimePeriod: "Nov 2023 – Jan 2025",
			Links: []pages.Link{
				"https://sakah.co/",
				"https://www.llamaevents.com/",
			},
			Description: []pages.Paragraph{
				"Built and shipped multiple Next.js applications including customer-facing platforms and event management systems.",
				"Led migration from Next.js 13 to 14, integrated Supabase replacing a Kubernetes/Python backend, and redesigned the database for improved scalability.",
			},
		},
	}

	_ = pages.ResumePage(experiences).Render(r.Context(), w)
}
