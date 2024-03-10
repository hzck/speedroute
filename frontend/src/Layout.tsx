import * as React from "react";
import { Box, Container } from "@mui/material";
import Copyright from "./Copyright";
import Logo from "./Logo";
import { Outlet } from "react-router-dom";

export default function Layout() {
  return (
    <Container maxWidth="sm">
      <Box sx={{ my: 4 }}>
        <Logo />
        <Outlet />
        <Copyright />
      </Box>
    </Container>
  );
}
