# Canids

## To Begin development

- To install dependencies, run:

```bash
npm i
```

- To run the development server:

```bash
npm start
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

- To build env:

```bash
npm run build
```

- To run build env:

```bash
npm run start:prod
```

- To run linter:

```bash
npm run lint
```

- To run prettier:

```bash
npm run format
```

- To run prettier fix:

```bash
npm run format:fix
```

## Important to Know

### To learn more about Next.js, take a look at the following resources:

- We are using Pages Router[Next.js Documentation](https://nextjs.org/docs/pages/building-your-application/routing/pages-and-layouts) - learn about Next.js features and API.

### For Styling, Grids and Icons we are using MUI Style Kit

- [Mui](https://mui.com/) - learn more about that.
- [Grid2](https://mui.com/material-ui/react-grid2/) we are using Grid V2 (try to keep it as usual)
- [Icons](https://mui.com/material-ui/material-icons/) - Icons listed here.
- [Fonts\*](https://fonts.google.com/) - Custom Fonts setUp files (exmple - Raleway):
  - `/src/pages/_document.tsx`
  - `src/styles/theme.ts`
  - `src/styles/global.css`

### Auth and Request processing

- `src/hooks/useRequest.tsx` - Stick to use this hook for any `src/api/*` calls
- `src/context/authContext.tsx` - Global App Auth processing using [Cookies](https://www.npmjs.com/package/react-cookie) and JWT Token

## File structure

- `/public/` - For aliases.
- `/src/api/` - All API calls should be located here (see exampleApi.ts).
- `/src/components/` - For atomic design structure (read more [Atomic Design](https://atomicdesign.bradfrost.com/chapter-2/) or [Whole documentation](https://atomicdesign.bradfrost.com)). Modals, Layouts and Forms should be located here, try to keep solutions as showed in LoginForm (we are using [RHF](https://react-hook-form.com/) + [YUP](https://www.npmjs.com/package/yup))
- `/src/constants/` - For project global constants (Routes, types, etc.).
- `/src/context/` - For React Global state management check Auth and Notifications contexts.
- `/src/hooks/` - For React hooks (useRequest hook located here).
- `/src/pages/` - For Next JS base routed pages [learn here](https://nextjs.org/docs/pages/building-your-application/routing/pages-and-layouts).
- `/src/styles/` - For global Styling solutions and MaterialUI theme configuration.
- `/src/utils/` - For global complex logic solution (take a look at usefull fucntions inside).
