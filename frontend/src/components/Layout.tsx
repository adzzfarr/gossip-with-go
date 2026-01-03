import { Box } from "@mui/material";
import ForumAppBar from "./AppBar";

interface LayoutProps {
    children: React.ReactNode;
}

export default function Layout({ children }: LayoutProps) {
    return (
        <Box
            sx={{
                display: "flex",
                flexDirection: "column",
                minHeight: "100vh"
            }}
        >
            <ForumAppBar />
            <Box  component="main" sx={{ flexGrow: 1 }}>
                {children}
            </Box>
        </Box>
    );
}