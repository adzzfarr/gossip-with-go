import { Grid } from "@mui/material";
import { TopicCard } from "./TopicCard";
import type { Topic } from "../types";

export default function TopicsList({ topics }: { topics: Topic[] }) {
    return (
        <Grid container spacing={3}>
            {topics.map(
                topic => (
                    <Grid
                        size={{
                            xs: 12,
                            sm: 6,
                            md: 4,
                        }}
                        key={topic.topicID}
                    >
                        <TopicCard topic={topic} />
                    </Grid>

                )
            )}
        </Grid>
    );
}