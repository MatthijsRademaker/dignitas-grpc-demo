# gRPC Auction House: slides

Slidev deck on the private `@dignitas/slidev-theme` (Azure Artifacts feed, see `.npmrc`).

```bash
npm install     # needs access to the dignitas-se feed
npm run dev     # http://localhost:3030, presenter view at /presenter
npm run export  # PDF
```

Edit [slides.md](./slides.md). Custom visuals live in [components/](./components) and read the brand
colour from `--se-color-primary`, so they follow the theme in light and dark mode.
