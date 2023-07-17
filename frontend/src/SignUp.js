import * as React from "react";
import Avatar from "@mui/material/Avatar";
import Button from "@mui/material/Button";
import CssBaseline from "@mui/material/CssBaseline";
import TextField from "@mui/material/TextField";
import Box from "@mui/material/Box";
import PersonAddAltIcon from "@mui/icons-material/PersonAddAlt";
import Typography from "@mui/material/Typography";
import Container from "@mui/material/Container";

export default function SignUp() {
  const [usernameError, setUsernameError] = React.useState(false);
  const [passwordError, setPasswordError] = React.useState(false);

  const handleSubmit = (event) => {
    event.preventDefault();

    const data = new FormData(event.currentTarget);
    const formObj = {
      username: data.get("username"),
      password: data.get("password"),
    };

    setUsernameError(false);
    setPasswordError(false);

    if (!/^[\w]{2,30}$/.test(formObj.username)) {
      setUsernameError(true);
    }
    if (formObj.password.length < 8) {
      setPasswordError(true);
    }

    console.log(formObj);
  };

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
          <PersonAddAltIcon />
        </Avatar>
        <Typography component="h1" variant="h5">
          Sign up
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
            helperText={
              usernameError ? "Valid characters are A-Z a-z 0-9 _" : ""
            }
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
            helperText={
              passwordError
                ? "The password need to have at least 8 characters"
                : ""
            }
          />
          <Button type="submit" fullWidth variant="contained" sx={{ mt: 3 }}>
            Sign up
          </Button>
        </Box>
      </Box>
    </Container>
  );
}
