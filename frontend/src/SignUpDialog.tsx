import * as React from "react";
import Dialog, { DialogProps } from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import SignUp from "./SignUp";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";

export default function SignUpDialog(props: DialogProps) {
  return (
    <Dialog open={props.open} onClose={props.onClose}>
      <DialogContent>
        <SignUp />
      </DialogContent>
      <DialogActions>
        <Box m="auto">
          <Button onClick={() => props.onClose} sx={{ mt: -4 }}>
            Cancel
          </Button>
        </Box>
      </DialogActions>
    </Dialog>
  );
}
