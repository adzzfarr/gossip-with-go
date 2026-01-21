import { Grid } from "@mui/material";
import type { Post } from "../types";
import { PostCard } from "./PostCard";

export default function PostsList({ posts }: { posts: Post[] }) {
    return (
        <Grid container spacing={3}>
            {posts.map(
                (post) => (
                    <Grid
                        size={{
                            xs: 12,
                            sm: 6,
                            md: 4,
                        }}
                        key={post.postID}
                    >
                            <PostCard post={post} />
                    </Grid>   
                )
            )}
        </Grid>
    );
}