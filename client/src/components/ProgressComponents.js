import { LinearProgress, CircularProgress, Box } from "@mui/material";
import "../stylesheets/progressComponents.css";

export function ProgressBar({ progress, message }) {
    return (
    <Box className="progress-box">
        <div className="progress-message"><b>{message}</b></div>
        <LinearProgress className="progress-bar" variant="determinate" value={progress} />
    </Box>
    );

}

export function LoadingSpinner({message}) {
    return (
    <Box className="progress-box">
        <div className="progress-message"><b>{message}</b></div>
        <CircularProgress />
    </Box>
    );
}