package io.paketo;

import java.io.IOException;

import jakarta.servlet.annotation.WebServlet;
import jakarta.servlet.http.HttpServlet;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

@WebServlet("/runtime")
public class RuntimeServlet extends HttpServlet {

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws IOException {
        System.getProperties().forEach((key, value) -> {
            try {
                resp.getOutputStream().println((String) key + "=" + (String) value);
            } catch (IOException e) {
                e.printStackTrace();
            }
        });
    }
}
