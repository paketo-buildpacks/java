package io.paketo.cki;

import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

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
