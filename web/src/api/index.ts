export interface Profile {
    id?: string,
    handle: string,
    avatar: string,
}

export const profile = {
    get: (token: string, id: string) => fetch(`/api/profile/${id}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
    }),
    create: (token: string, profile: Profile) => fetch("/api/profile", {
        method: "POST",
        headers: {
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(profile),
    }),
    update: (token: string, profile: Profile) => fetch(`/api/profile`, {
        method: "PUT",
        headers: {
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(profile),
    }),
    delete: (token: string, id: string) => fetch(`/api/profile/${id}`, {
        method: "DELETE",
        headers: {
          'Authorization': `Bearer ${token}`,
        },
    }),
};

export interface Post {
    id?: number,
    author_id: string,
    content: string,
    type: string,
    data: string,
    visibility: string,
}

export const posts = {
    get: (token: string, id: string) => fetch(`/api/posts/${id}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        },
    }),
    list: (token: string, limit: number, offset: number) => fetch(`/api/posts?limit=${limit}&offset=${offset}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
        },
    }),
    create: (token: string, post: Post) => fetch("/api/posts", {
        method: "POST",
        headers: {
            'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(post),
    }),
    update: (token: string, post: Post) => fetch("/api/posts", {
        method: "PUT",
        headers: {
            'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(post),
    }),
    delete: (token: string, id: string) => fetch(`/api/posts/${id}`, {
        method: "DELETE",
        headers: {
            'Authorization': `Bearer ${token}`,
        },
    }),
};