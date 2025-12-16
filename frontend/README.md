# Suipic Frontend

Frontend application for Suipic photo management system built with SvelteKit.

## Tech Stack

- **Framework**: SvelteKit
- **Language**: TypeScript
- **Styling**: TailwindCSS + DaisyUI
- **Icons**: Iconify
- **State Management**: Svelte Stores
- **Data Fetching**: TanStack Query (Svelte Query)
- **Build Tool**: Vite
- **Package Manager**: pnpm

## Getting Started

### Prerequisites

- Node.js 18+
- pnpm 8+

### Installation

```bash
pnpm install
```

### Development

```bash
pnpm dev
```

The app will be available at `http://localhost:5173`

### Build

```bash
pnpm build
```

### Preview Production Build

```bash
pnpm preview
```

## Project Structure

```
src/
├── lib/
│   ├── api/           # API client and endpoints
│   ├── components/    # Reusable Svelte components
│   ├── stores/        # Svelte stores for state management
│   ├── types/         # TypeScript type definitions
│   └── utils/         # Utility functions
├── routes/            # SvelteKit routes (pages)
│   ├── albums/        # Album-related pages
│   ├── login/         # Login page
│   ├── register/      # Registration page
│   ├── profile/       # User profile page
│   ├── settings/      # Settings page
│   └── +layout.svelte # Root layout
├── app.css            # Global styles
└── app.html           # HTML template
```

## Features

- **Authentication**: JWT-based auth with token storage and auto-refresh
- **Albums Management**: 
  - Create, edit, and delete photo albums
  - Grid and table view for album listing
  - Custom fields for flexible metadata
  - Client assignment for sharing albums
  - Thumbnail selection from album photos
  - Album detail view with photo gallery
- **Photo Upload**: Upload photos to albums with drag-and-drop support
- **Responsive Design**: Mobile-friendly UI with responsive navigation
- **Dark Mode**: Theme switching between light and dark modes
- **Type Safety**: Full TypeScript support with backend DTO matching

See [ALBUM_MANAGEMENT.md](./ALBUM_MANAGEMENT.md) for detailed documentation on album management features.

## API Integration

The frontend connects to the backend API at `http://localhost:8080` (proxied through Vite).

All API calls include JWT authentication via the Authorization header when required.

## Conventions

- Types are prefixed with `T` (e.g., `TUser`, `TPhoto`)
- Enums are prefixed with `E` (e.g., `EUserRole`, `EPickRejectState`)
- All types are in separate `*.types.ts` files
- Stores are centralized in `src/lib/stores/`
- API clients are modular in `src/lib/api/`
