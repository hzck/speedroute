import * as React from "react";
import { CssBaseline, ThemeProvider } from "@mui/material";
import theme from "./theme";

import { RouterProvider, createBrowserRouter } from "react-router-dom";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout />,
    loader: () => ({ message: "Hello TEST!" }),
    children: [
      {
        path: "login",
        element: <Login />,
      },
      {
        path: "settings",
        element: <Settings />,
      },
      {
        path: "/:id",
        element: <div>id</div>,
      },
    ],
  },
]);

export default function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <RouterProvider router={router} fallbackElement={<p>Loading...</p>} />;
    </ThemeProvider>
  );
}
