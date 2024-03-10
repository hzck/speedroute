import * as React from "react";
import Avatar from "@mui/material/Avatar";
import Button from "@mui/material/Button";
import CssBaseline from "@mui/material/CssBaseline";
import TextField from "@mui/material/TextField";
import Box from "@mui/material/Box";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import Typography from "@mui/material/Typography";
import Container from "@mui/material/Container";
import { Divider } from "@mui/material";
import SignUpDialog from "./SignUpDialog";
import { texts } from "./texts/texts";

export default function Login() {
  const [usernameError, setUsernameError] = React.useState(false);
  const [passwordError, setPasswordError] = React.useState(false);
  const [open, setOpen] = React.useState(false);

  const handleSubmit = (event) => {
    event.preventDefault();

    const data = new FormData(event.currentTarget);
    const formObj = {
      username: data.get("username"),
      password: data.get("password"),
    };

    setUsernameError(false);
    setPasswordError(false);

    if (!/^[\w]{2,30}$/.test(formObj.username as string)) {
      setUsernameError(true);
    }
    if ((formObj.password as string).length < 8) {
      setPasswordError(true);
    }

    console.log(formObj);
  };

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <Box
        sx={{
          marginBottom: 8,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
          <LockOutlinedIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Login
        </Typography>
        <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            id="username"
            label="Username"
            name="username"
            autoComplete="username"
            error={usernameError}
            helperText={usernameError ? texts.login.userNameTaken : ""}
            autoFocus
          />
          <TextField
            margin="normal"
            required
            fullWidth
            name="password"
            label="Password"
            type="password"
            id="password"
            autoComplete="current-password"
            error={passwordError}
            helperText={passwordError ? texts.login.passwordHelper : ""}
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 3 }}
          >
            Login
          </Button>
        </Box>
        <Divider role="presentation" flexItem>
          OR
        </Divider>
        <Button variant="outlined" onClick={handleClickOpen} sx={{ mt: 3 }}>
          Sign up
        </Button>
      </Box>
      <SignUpDialog open={open} onClose={handleClose} />
    </Container>
  );
}
